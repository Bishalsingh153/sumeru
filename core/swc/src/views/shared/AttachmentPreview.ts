import { html } from "../../template/html.js";

export interface AttachmentPreviewTarget {
  name: string;
  url: string;
  mimetype?: string;
}

/** Prefer inline display for PDF/image in preview iframes and img tags. */
export function attachmentContentInlineUrl(url: string): string {
  if (!url) return url;
  if (url.includes("download=")) return url;
  const sep = url.includes("?") ? "&" : "?";
  return `${url}${sep}download=0`;
}

export function attachmentPreviewModal(
  target: AttachmentPreviewTarget | null,
  onClose: () => void,
): ReturnType<typeof html> {
  if (!target) return html``;
  const mime = (target.mimetype ?? "").toLowerCase();
  const isPDF = mime === "application/pdf" || target.name.toLowerCase().endsWith(".pdf");
  const isImage = mime.startsWith("image/");
  const inlineUrl = attachmentContentInlineUrl(target.url);
  return html`
    <div class="sum-attachment-preview-backdrop" @click=${onClose}>
      <div class="sum-attachment-preview" @click=${(e: Event) => e.stopPropagation()}>
        <header class="sum-attachment-preview-header">
          <span>${target.name}</span>
          <button type="button" class="sum-btn sum-btn--ghost" @click=${onClose}>Close</button>
        </header>
        <div class="sum-attachment-preview-body">
          ${isImage
            ? html`<img class="sum-attachment-preview-img" src=${inlineUrl} alt=${target.name} />`
            : isPDF
              ? html`<iframe
                  class="sum-attachment-preview-pdf"
                  src=${inlineUrl}
                  title=${target.name}
                  sandbox="allow-same-origin"
                ></iframe>
                <p class="sum-msg-form-hint">
                  <a class="sum-chatter-attachment-link" href=${inlineUrl} target="_blank" rel="noopener">Open PDF in new tab</a>
                </p>`
              : html`<a class="sum-field-link" href=${target.url} download>Download ${target.name}</a>`}
        </div>
      </div>
    </div>
  `;
}
