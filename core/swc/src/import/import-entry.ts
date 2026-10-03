import { mountImportWizard } from "./ImportWizard.js";

const root = document.getElementById("sum-import-root");
if (root) {
  const batch = root.dataset.batch ?? "";
  const csrf = root.dataset.csrf ?? "";
  mountImportWizard(document.getElementById("sum-import-app"), batch, csrf);
}
