import { SwcComponent } from "../runtime/component.js";
import { html } from "../template/html.js";
import { SWC_API_BASE } from "../constants/routes.js";
import { inputValueFromEvent } from "../widgets/field-events.js";

interface DirectUser {
  id: number;
  name: string;
  login: string;
}

interface DirectConversation {
  userId: number;
  name: string;
  preview?: string;
  lastMessage?: string;
}

interface DirectAttachment {
  id: number;
  name: string;
  url: string;
}

interface DirectMessage {
  id: number;
  body: string;
  createDate: string;
  outgoing: boolean;
  attachments?: DirectAttachment[];
}

export class InternalMessagesPanel extends SwcComponent {
  private search = "";
  private searchResults: DirectUser[] = [];
  private conversations: DirectConversation[] = [];
  private peerId = 0;
  private peerName = "";
  private thread: DirectMessage[] = [];
  private draft = "";
  private loading = false;
  private posting = false;
  private enabled = true;
  private disabledReason = "";
  private searchTimer: ReturnType<typeof setTimeout> | null = null;

  override setup(): void {
    void this.refreshConversations();
  }

  private apiBase(): string {
    return this.env.bootstrap.swcApiBase || SWC_API_BASE;
  }

  private async refreshConversations(): Promise<void> {
    try {
      const data = await this.env.services.http.getJSON<{ conversations: DirectConversation[] }>(
        `${this.apiBase()}/direct/conversations`,
      );
      this.conversations = data.conversations ?? [];
      this.enabled = true;
      this.rerender();
    } catch {
      this.enabled = false;
      this.disabledReason = "Internal messages are unavailable or disabled for this company.";
      this.rerender();
    }
  }

  private async runSearch(): Promise<void> {
    const q = this.search.trim();
    if (q.length < 1) {
      this.searchResults = [];
      this.rerender();
      return;
    }
    try {
      const data = await this.env.services.http.getJSON<{ users: DirectUser[] }>(
        `${this.apiBase()}/direct/users?q=${encodeURIComponent(q)}`,
      );
      this.searchResults = data.users ?? [];
    } catch {
      this.searchResults = [];
    }
    this.rerender();
  }

  private scheduleSearch(): void {
    if (this.searchTimer) clearTimeout(this.searchTimer);
    this.searchTimer = setTimeout(() => void this.runSearch(), 250);
  }

  private async openPeer(userId: number, name: string): Promise<void> {
    this.peerId = userId;
    this.peerName = name;
    this.search = "";
    this.searchResults = [];
    await this.loadThread();
  }

  private async loadThread(): Promise<void> {
    if (this.peerId <= 0) return;
    this.loading = true;
    this.rerender();
    try {
      const data = await this.env.services.http.getJSON<{ messages: DirectMessage[] }>(
        `${this.apiBase()}/direct/thread?userId=${this.peerId}`,
      );
      this.thread = data.messages ?? [];
    } finally {
      this.loading = false;
      this.rerender();
    }
  }

  private async postMessage(): Promise<void> {
    const body = this.draft.trim();
    if (!body || this.peerId <= 0) return;
    this.posting = true;
    this.rerender();
    try {
      await this.env.services.http.postForm(`${this.apiBase()}/direct/post`, {
        recipient_id: String(this.peerId),
        body,
      });
      this.draft = "";
      await this.loadThread();
      await this.refreshConversations();
    } finally {
      this.posting = false;
      this.rerender();
    }
  }

  private async uploadFile(file: File): Promise<void> {
    if (this.peerId <= 0) return;
    let messageId = 0;
    if (this.draft.trim()) {
      const res = await this.env.services.http.postForm(`${this.apiBase()}/direct/post`, {
        recipient_id: String(this.peerId),
        body: this.draft.trim(),
      });
      const parsed = (await res.json()) as { id?: number };
      messageId = parsed.id ?? 0;
      this.draft = "";
    } else if (this.thread.length > 0) {
      messageId = this.thread[this.thread.length - 1].id;
    } else {
      const res = await this.env.services.http.postForm(`${this.apiBase()}/direct/post`, {
        recipient_id: String(this.peerId),
        body: " ",
      });
      const parsed = (await res.json()) as { id?: number };
      messageId = parsed.id ?? 0;
    }
    if (messageId <= 0) return;
    const form = new FormData();
    form.set("message_id", String(messageId));
    form.set("file", file);
    await this.env.services.http.postMultipart("/web/swc/direct/upload", form);
    await this.loadThread();
    await this.refreshConversations();
    this.rerender();
  }

  override template() {
    if (!this.enabled) {
      return html`<div class="sum-msg-empty" role="status">
        <p class="sum-msg-empty-title">Messages</p>
        <p class="sum-msg-empty-hint">${this.disabledReason}</p>
      </div>`;
    }

    const listItems =
      this.searchResults.length > 0
        ? this.searchResults.map(
            (u) => html`<li>
              <button type="button" class="sum-msg-peer-btn" @click=${() => void this.openPeer(u.id, u.name || u.login)}>
                <span class="sum-msg-author">${u.name || u.login}</span>
                <span class="sum-msg-time">${u.login}</span>
              </button>
            </li>`,
          )
        : this.conversations.map(
            (c) => html`<li>
              <button
                type="button"
                class="sum-msg-peer-btn${this.peerId === c.userId ? " is-active" : ""}"
                @click=${() => void this.openPeer(c.userId, c.name)}
              >
                <span class="sum-msg-author">${c.name}</span>
                <span class="sum-msg-time">${c.preview ?? ""}</span>
              </button>
            </li>`,
          );

    return html`<div class="sum-msg-shell sum-msg-shell--dock">
      <div class="sum-msg-layout">
        <div class="sum-msg-sidebar">
          <input
            type="search"
            class="sum-msg-search"
            placeholder="Search internal users…"
            value=${this.search}
            @input=${(event: Event) => {
              this.search = inputValueFromEvent(event);
              this.scheduleSearch();
              this.rerender();
            }}
          />
          <ul class="sum-msg-peer-list">${listItems}</ul>
        </div>
        <div class="sum-msg-main">
          ${this.peerId <= 0
            ? html`<div class="sum-msg-thread-empty">
                <p class="sum-msg-thread-empty-title">Direct messages</p>
                <p class="sum-msg-thread-empty-hint">Search for a colleague or pick a recent conversation.</p>
              </div>`
            : html`<div class="sum-msg-thread-head">${this.peerName}</div>
              <div class="sum-msg-thread">
                ${this.loading
                  ? html`<p class="sum-msg-thread-empty-hint">Loading…</p>`
                  : this.thread.length === 0
                    ? html`<p class="sum-msg-thread-empty-hint">No messages yet. Say hello.</p>`
                    : this.thread.map(
                        (m) => html`<div class="sum-msg-card${m.outgoing ? " sum-msg-card--out" : ""}">
                          <div class="sum-msg-card-head">
                            <span class="sum-msg-time">${m.createDate}</span>
                          </div>
                          <div class="sum-msg-body">${m.body}</div>
                          ${(m.attachments ?? []).length > 0
                            ? html`<ul class="sum-chatter-attachments--inline">
                                ${(m.attachments ?? []).map(
                                  (a) => html`<li><a class="sum-chatter-attachment-link" href=${a.url}>${a.name}</a></li>`,
                                )}
                              </ul>`
                            : ""}
                        </div>`,
                      )}
              </div>
              <div class="sum-msg-composer">
                <div class="sum-msg-form">
                  <textarea
                    class="sum-msg-input"
                    rows="3"
                    placeholder="Message ${this.peerName}…"
                    value=${this.draft}
                    @input=${(event: Event) => {
                      this.draft = inputValueFromEvent(event);
                      this.rerender();
                    }}
                  ></textarea>
                  <div class="sum-msg-form-actions">
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
                      <span class="sum-btn sum-btn--secondary">Attach</span>
                    </label>
                    <button
                      type="button"
                      class="sum-msg-send"
                      disabled=${this.posting ? "disabled" : undefined}
                      @click=${() => void this.postMessage()}
                    >
                      Send
                    </button>
                  </div>
                  <p class="sum-msg-form-hint">Internal users only · not shown on record Log</p>
                </div>
              </div>`}
        </div>
      </div>
    </div>`;
  }
}
