package module

import "context"

// ExpandInstallModuleNamesForTest exposes expandInstallModuleNames for tests.
func ExpandInstallModuleNamesForTest(ctx context.Context, parts []string) ([]string, error) {
	return expandInstallModuleNames(ctx, parts)
}

func ResolveModuleCategoryIDForTest(ctx context.Context, categoryName string) (interface{}, error) {
	return resolveModuleCategoryID(ctx, categoryName)
}

const (
	ModuleReloadInstallForTest = moduleReloadInstall
	ModuleReloadUpdateForTest  = moduleReloadUpdate
)

type DataFileOptsForTest = dataFileOpts

func NewDataFileOptsForTest(noUpdate bool) DataFileOptsForTest {
	return dataFileOpts{noUpdate: noUpdate}
}

func (o DataFileOptsForTest) SkipExistingOnUpdateForTest(ctx context.Context, moduleName, xmlID string) bool {
	return dataFileOpts(o).skipExistingOnUpdate(ctx, moduleName, xmlID)
}

// RegistryHasFieldForTest reports whether a model field exists in the ORM registry.
func RegistryHasFieldForTest(modelName, fieldName string) bool {
	return registryHasField(modelName, fieldName)
}

// FieldAccessDefaultNameForTest builds the default sys.field.access name for CSV rows.
func FieldAccessDefaultNameForTest(model, field, groupXML string) string {
	return fieldAccessDefaultName(model, field, groupXML)
}
