package image

import (
	"testing"
	"minik8s/apiobjects"
)


func TestGenetateDockerFile(t *testing.T) {
    input := apiobjects.FunctionCtlInput{
		Object: apiobjects.Object{
			TypeMeta: apiobjects.TypeMeta{
			   ApiVersion: "app/v1",
			   Kind:       "Function",
			},
			ObjectMeta: apiobjects.ObjectMeta{
				Name: "dns-function",
			},
		},
		FunctionSpec: apiobjects.FunctionSpec{
			MinReplicas:         1,
			MaxReplicas:         2,
			TargetQPSPerReplica: 100,
		},
		BuildOptions: apiobjects.BuildOptions{
			ExtraCommands:   []string{"RUN apt-get update","RUN apt-get install -y curl"},
			FunctionFileDir: "/test",
		},
	}
	err := GenerateDockerfile(input)
	if err != nil {
		t.Error("GenerateDockerfile error")
	}
}
