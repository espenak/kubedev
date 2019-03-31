package main

type TemplateVariables struct {
	context Context
}

func (templateVariables *TemplateVariables) FullName() string {
	return templateVariables.context.FullName()
}

func (templateVariables *TemplateVariables) RootDirectory() string {
	return templateVariables.context.RootDirectory
}

func (templateVariables *TemplateVariables) Path(key string) string {
	return templateVariables.context.Paths[key]
}

func (templateVariables *TemplateVariables) Var(key string) string {
	return templateVariables.context.Vars[key]
}

func (templateVariables *TemplateVariables) UserConfig(key string) string {
	return templateVariables.context.UserConfig[key]
}
