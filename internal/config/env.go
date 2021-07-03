package config

import (
	"os"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// EnvConfig
type EnvConfig struct {
	Cluster         string `export:"EKS_CLUSTER_NAME" default:"development-lendi-eks-cluster"`
	Region          string `export:"AWS_REGION" default:"ap-southeast-2"`
	JanitorLabel    string `export:"JANITOR_LABEL" default:"helm-janitor=true"`
	TmpFileLocation string `export:"TMP_FILE_LOCATION" default:"/tmp"`
	TmpFilePrefix   string `export:"TMP_FILE_PREFIX" default:"k8s-ca-*"`
	DebugFlag       bool   `export:"DEBUG" default:"false"`
}

func (c *EnvConfig) Init() {
	typ := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f, _ := typ.FieldByName(typeOfS.Field(i).Name)
		envkey := f.Tag.Get("export")
		if os.Getenv(envkey) != "" {
			// c[typeOfS.Field(i).Name] = os.Getenv(envkey)
			// reflect.ValueOf(c).Elem().FieldByName(f.Name).SetString(os.Getenv(envkey))
			c.setFromEnv(f.Name)
		} else {
			c.setDefaultEnv(f.Name)
			// reflect.ValueOf(c).Elem().FieldByName(f.Name).SetString(def)
		}
		log.Debugf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
}

// Populates the EnvConfig struct with...
func (c *EnvConfig) setFromEnv(field string) {
	typ := reflect.TypeOf(*c)
	f, found := typ.FieldByName(field)
	if !found {
		log.Fatalf("Received an invalid struct field: %s which is not part of EnvConfig", field)
	}
	envkey := f.Tag.Get("export")
	envVal := os.Getenv(envkey)
	if f.Type.Name() == "bool" {
		a, _ := strconv.ParseBool(envVal)
		reflect.ValueOf(c).Elem().FieldByName(f.Name).SetBool(a)
	}
	if f.Type.Name() == "string" {
		reflect.ValueOf(c).Elem().FieldByName(f.Name).SetString(envVal)
	}
}

// Populates the EnvConfig struct with default value depending on the field
func (c *EnvConfig) setDefaultEnv(field string) {
	typ := reflect.TypeOf(*c)
	f, found := typ.FieldByName(field)
	if !found {
		log.Fatalf("Received an invalid struct field: %s which is not part of EnvConfig", field)
	}
	def := f.Tag.Get("default")
	if f.Type.Name() == "bool" {
		a, _ := strconv.ParseBool(def)
		reflect.ValueOf(c).Elem().FieldByName(f.Name).SetBool(a)
	}
	if f.Type.Name() == "string" {
		reflect.ValueOf(c).Elem().FieldByName(f.Name).SetString(def)
	}
}
