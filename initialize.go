package goone

import (
	"golang.org/x/tools/go/analysis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	Analyzer.Flags.StringVar(&configPath, "configPath", "", "config file path(abs)")
}

func appendTypes(pass *analysis.Pass, pkg, name string) {

	//TODO Use types instead of the name of types
	//if typ := analysisutil.TypeOf(pass, pkg, name); typ != nil {
	//	sqlTypes = append(sqlTypes, typ)
	//} else {
	//	if name[0] == '*' {
	//		name = name[1:]
	//		pkg = "*" + pkg
	//	}
	//	for _, v := range pass.TypesInfo.Types {
	//		if v.Type.String() == pkg+"."+name {
	//			sqlTypes = append(sqlTypes, v.Type)
	//			return
	//		}
	//	}
	//}

	if name[0] == '*' {
		name = name[1:]
		pkg = "*" + pkg
	}
	sqlTypes = append(sqlTypes, pkg+"."+name)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func readTypeConfig() *Types {
	var cp string
	// if configPath flag is not set
	if configPath ==""{

		curDir, _ := os.Getwd()
		for !fileExists(curDir+"/goone.yml"){
			// Search up to the root
			if curDir == filepath.Dir(curDir) || curDir == ""{
				// If goone.yml is not found
				return nil
			}
			curDir = filepath.Dir(curDir)
		}
		cp = curDir+"/goone.yml"
	}else{
		cp = configPath
	}

	buf, err := ioutil.ReadFile(cp)
	if err != nil {
		return nil
	}
	typesFromConfig := Types{}
	err = yaml.Unmarshal([]byte(buf), &typesFromConfig)
	if err != nil {
		log.Fatalf("yml parse error:%v", err)
	}

	return &typesFromConfig
}

func prepareTypes(pass *analysis.Pass){
	typesFromConfig := readTypeConfig()
	if typesFromConfig != nil {
		for _, pkg := range typesFromConfig.Package {
			pkgName := pkg.PkgName
			for _, typeName := range pkg.TypeNames {
				appendTypes(pass, pkgName, typeName.TypeName)
			}
		}
	}

	appendTypes(pass, "database/sql", "*DB")
	appendTypes(pass, "gorm.io/gorm", "*DB")
	appendTypes(pass, "gopkg.in/gorp.v1", "*DbMap")
	appendTypes(pass, "gopkg.in/gorp.v2", "*DbMap")
	appendTypes(pass, "github.com/go-gorp/gorp/v3", "*DbMap")
	appendTypes(pass, "github.com/jmoiron/sqlx", "*DB")
}