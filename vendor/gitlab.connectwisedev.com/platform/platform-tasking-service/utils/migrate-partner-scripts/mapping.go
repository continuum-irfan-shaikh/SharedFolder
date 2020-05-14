package main

var templateToOrigin = map[int]string{
	bashTemplateIDint:       "37f7f19f-40e8-11e9-a643-e0d55e1ce78a", //TODO: according to https://gitlab.connectwisedev.com/platform/rmm-scripts/pull/196
	cmdTemplateIDint:        "e3d2c26b-c5ba-49cf-a089-7637f6de949e",
	powershellTemplateIDint: "51a74346-e19b-11e7-9809-0800279505d9",
}

var converters = map[int]func([]ScriptMsSQL) ([]TaskDefinition, error){
	cmdTemplateIDint:        cmdDefs,
	powershellTemplateIDint: getPowerShellTaskDef,
	//bashTemplateIDint:       getBashTaskDef,
}
