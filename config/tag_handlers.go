package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
)

var loadSecrets = sync.OnceValue(func() map[string]string {
	secretsPath := filepath.Join(utils.ExeDirectory, "secrets.yaml")
	data, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", Path()).Msg("fatal error reading \"secrets.yaml\"")
	}
	var secrets map[string]string
	if err = yaml.Unmarshal(data, &secrets); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing \"secrets.yaml\"")
	}
	return secrets
})

func mustGetSecret(key string) string {
	secrets := loadSecrets()
	if secret, ok := secrets[key]; ok {
		return secret
	}
	log.Fatal().Msgf("secret \"%s\" not found in \"secrets.yaml\"", key)
	return ""
}

type tagVisitor struct{}

func (v *tagVisitor) Visit(node ast.Node) ast.Visitor {
	if mappingNode, ok := node.(*ast.MappingValueNode); ok {
		if tagNode, ok := mappingNode.Value.(*ast.TagNode); ok {
			if stringNode, ok := tagNode.Value.(*ast.StringNode); ok {
				tagName := tagNode.GetToken().Value
				if tagName == "!include" {
					includeFilePath := stringNode.Value
					if !path.IsAbs(includeFilePath) {
						includeFilePath = filepath.Join(filepath.Dir(Path()), includeFilePath)
					}

					parsedAst, err := parser.ParseFile(includeFilePath, parser.ParseComments)
					if err != nil {
						log.Fatal().Err(err).Str("path", includeFilePath).Msg("fatal error loading file")
					}

					mappingNode.Value = parsedAst.Docs[0].Body
				}
				if tagName == "!secret" {
					stringNode.Value = mustGetSecret(stringNode.Value)
					mappingNode.Value = stringNode
				}

			}
		}
	}
	return v
}

func processTags(bytes []byte) (ast.Node, error) {
	parsedAst, err := parser.ParseBytes(bytes, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(parsedAst.Docs) == 0 {
		return nil, fmt.Errorf("no documents found in config")
	}
	rootNode := parsedAst.Docs[0].Body
	ast.Walk(&tagVisitor{}, rootNode)
	return rootNode, nil
}
