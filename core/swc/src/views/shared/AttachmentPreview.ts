import { html } from "../../template/html.js";

export interface AttachmentPreviewTarget {
  name: string;
  url: string;
  mimetype?: string;
}

export function attachmentPreviewModal(
  target: AttachmentPreviewTarget | null,
  onClose: () => void,
): ReturnType<typeof html> {
  if (!target) return html``;
  const mime = (target.mimetype ?? "").toLowerCase();
  const isPDF = mime === "application/pdf" || target.name.toLowerCase().endsWith(".pdf");
  const isImage = mime.startsWith("image/");
  return html`
    <div class="sum-attachment-preview-backdrop" @click=${onClose}>
      <div class="sum-attachment-preview" @click=${(e: Event) => e.stopPropagation()}>
        <header class="sum-attachment-preview-header">
          <span>${target.name}</span>
          <button type="button" class="sum-btn sum-btn--ghost" @click=${onClose}>Close</button>
        </header>
        <div class="sum-attachment-preview-body">
          ${isImage
            ? html`<img class="sum-attachment-preview-img" src=${target.url} alt=${target.name} />`
            : isPDF
              ? html`<iframe class="sum-attachment-preview-pdf" src=${target.url} title=${target.name} sandbox=""></iframe>`
              : html`<a class="sum-field-link" href=${target.url} download>Download ${target.name}</a>`}
        </div>
      </div>
    </div>
  `;
}
