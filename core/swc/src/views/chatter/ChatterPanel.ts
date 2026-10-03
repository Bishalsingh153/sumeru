import { SwcComponent } from "../../runtime/component.js";
import { html } from "../../template/html.js";
import { RECORD_UPDATED, SWC_API_BASE } from "../../constants/routes.js";
import { inputValueFromEvent } from "../../widgets/field-events.js";
import { attachmentPreviewModal, type AttachmentPreviewTarget } from "../shared/AttachmentPreview.js";

interface ChatterMessage {
  body: string;
  author: string;
  createDate: string;
}

interface ChatterAttachment {
  id: number;
  name: string;
  url: string;
  mimetype?: string;
}

interface ChatterPayload {
  model: string;
  recordId: number;
  messages: ChatterMessage[];
  attachments?: ChatterAttachment[];
  enabled: boolean;
}

interface ChatterPanelProps {
  model: string;
  recordId: number;
  csrfToken: string;
  /** When true, omit outer aside chrome (activity dock mount). */
  embedded?: boolean;
}

export class ChatterPanel extends SwcComponent<ChatterPanelProps> {
  private messages: ChatterMessage[] = [];
  private attachments: ChatterAttachment[] = [];
  private draft = "";
  private loading = true;
  private posting = false;
  private enabled = true;
  private uploadPct = 0;
  private uploading = false;
  private preview: AttachmentPreviewTarget | null = null;
  private unsubRecordUpdated: (() => void) | null = null;

  override setup(): void {
    void this.load();
  }

  override onMount(): void {
    this.unsubRecordUpdated = this.env.services.bus.subscribe(RECORD_UPDATED, (payload) => {
      const msg = payload as { model?: string; id?: number; recordId?: number };
      const rid = msg.id ?? msg.recordId;
      if (msg.model !== this.props.model) return;
      if (rid != null && rid > 0 && rid !== this.props.recordId) return;
      void this.load();
    });
  }

  override onWillUnmount(): void {
    this.unsubRecordUpdated?.();
    this.unsubRecordUpdated = null;
  }

  override onPropsChanged(props: ChatterPanelProps): void {
    if (props.model !== this.props.model || props.recordId !== this.props.recordId) {
      void this.load();
    }
  }

  private async load(): Promise<void> {
    const { model, recordId } = this.props;
    if (recordId <= 0) {
      this.loading = false;
      this.rerender();
      return;
    }
    this.loading = true;
    this.rerender();
    try {
      const base = this.env.bootstrap.swcApiBase || SWC_API_BASE;
      const data = await this.env.services.http.getJSON<ChatterPayload>(
        `${base}/chatter?model=${encodeURIComponent(model)}&id=${recordId}`,
      );
      this.messages = data.messages ?? [];
      this.attachments = data.attachments ?? [];
      this.enabled = data.enabled !== false;
    } finally {
      this.loading = false;
      this.rerender();
    }
  }

  private async uploadFile(file: File): Promise<void> {
    if (this.props.recordId <= 0) return;
    this.uploading = true;
    this.uploadPct = 0;
    this.rerender();
    try {
      const form = new FormData();
      form.set("model", this.props.model);
      form.set("res_id", String(this.props.recordId));
      form.set("file", file);
      const res = await this.env.services.http.postMultipart("/web/chatter/upload", form, (pct) => {
        this.uploadPct = pct;
        this.rerender();
      });
      if (!res.ok) {
        return;
      }
      await this.load();
    } finally {
      this.uploading = false;
      this.uploadPct = 0;
      this.rerender();
    }
  }

  private async post(): Promise<void> {
    const body = this.draft.trim();
    if (!body || this.props.recordId <= 0) return;
    this.posting = true;
    this.rerender();
    try {
      await this.env.services.http.postForm("/web/chatter/post", {
        model: this.props.model,
        res_id: String(this.props.recordId),
        body,
        next: window.location.pathname + window.location.search,
      });
      this.draft = "";
      await this.load();
      this.env.services.bus.emit(RECORD_UPDATED, {
        model: this.props.model,
        id: this.props.recordId,
        recordId: this.props.recordId,
      });
    } finally {
      this.posting = false;
      this.rerender();
    }
  }

  override template() {
    if (this.props.recordId <= 0) {
      return html`<div class="sum-chatter sum-chatter--empty">Save the record to post messages.</div>`;
    }
    if (this.loading) {
      return html`<div class="sum-chatter sum-chatter--loading">Loading messages…</div>`;
    }
    if (!this.enabled) {
      return html`<div class="sum-chatter sum-chatter--disabled">
        <p>Internal messages are disabled for this company. Enable chatter on the company record to post comments.</p>
      </div>`;
    }
    const inner = html`
      <div class="sum-chatter${this.props.embedded ? " sum-chatter--embedded" : ""}">
        <ul class="sum-chatter-messages">
          ${this.messages.length === 0
            ? html`<li class="sum-chatter-empty">No messages yet.</li>`
            : this.messages.map(
                (m) => html`<li class="sum-chatter-message">
                  <div class="sum-chatter-meta">${m.author} · ${m.createDate}</div>
                  <div class="sum-chatter-body">${m.body}</div>
                </li>`,
              )}
        </ul>
        ${this.attachments.length > 0
          ? html`<ul class="sum-chatter-attachments sum-chatter-attachments--inline">
              ${this.attachments.map(
                (a) => html`<li>
                  <button
                    type="button"
                    class="sum-chatter-attachment-link"
                    @click=${() => {
                      this.preview = { name: a.name, url: a.url, mimetype: a.mimetype };
                      this.rerender();
                    }}
                  >
                    ${a.name}
                  </button>
                </li>`,
              )}
            </ul>`
          : ""}
        <div class="sum-chatter-composer">
          <textarea
            class="sum-chatter-input"
            placeholder="Write an internal message…"
            rows="3"
            value=${this.draft}
            @input=${(event: Event) => {
              this.draft = inputValueFromEvent(event);
              this.rerender();
            }}
          ></textarea>
          <div class="sum-chatter-composer-actions">
            <label class="sum-chatter-upload">
              <input
                type="file"
                hidden
                @change=${(event: Event) => {
                  const input = event.target as HTMLInputElement;
                  const file = input.files?.[0];
                  input.value = "";
                  if (file) void this.uploadFile(file);
                }}
              />
              <span class="sum-btn sum-btn--secondary">${this.uploading ? `Uploading ${this.uploadPct}%` : "Attach"}</span>
            </label>
            <button
              type="button"
              class="sum-btn sum-btn--primary sum-chatter-send"
              disabled=${this.posting ? "disabled" : undefined}
              @click=${() => void this.post()}
            >
              Post
            </button>
          </div>
        </div>
        ${attachmentPreviewModal(this.preview, () => {
          this.preview = null;
          this.rerender();
        })}
      </div>
    `;
    if (this.props.embedded) {
      return html`<div class="sum-msg-shell">${inner}</div>`;
    }
    return html`<aside class="sum-chatter-host">${inner}</aside>`;
  }
}
