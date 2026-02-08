package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}
type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"` //from yaml||env import env variable, from cleanenv
	/*“This Go field is called Env,
	but in YAML it’s named env,
	in the OS it’s called ENV,a
	it must exist,
	and if not, default to production.”*/
	Storage_path string     `yaml:"storage_path" env-required:"true"`
	HTTPServer   HTTPServer `yaml:"http_server"` //these are called stuct tags
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		//flags are variable within the terminal request curl skfj/4544 -name oki
		flags := flag.String("config", "", "path to the config file") //String(name string, value string, usage string) *string,here value is the deafult value and usage is the statement showed when someone asks what is this in help

		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("configpath is not set")
		}
	}
	fmt.Println(configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) { //Stat(name string) (os.FileInfo, error) , Stat returns a [FileInfo] describing the named file. If there is an error, it will be of type [*PathError].
		log.Fatal("config file does not exist")
	}

	var cfg Config
	fmt.Println(configPath)
	err := cleanenv.ReadConfig(configPath, &cfg) //It reads your config file (YAML / ENV) from configPath and fills the cfg struct with the values.
	if err != nil {
		log.Fatal("couldnt pass the values to cfg struct")

	}

	return &cfg
}

/*
if the file at configPath does NOT exist {
    // run this block
}

2. _, err := os.Stat(configPath)
_ means: ignore the fileInfo result
err gets the error (if any)
So now:
err == nil → file exists
err != nil → something went wrong
3. os.IsNotExist(err)
This is a helper function that checks:
“Is this error specifically because the file does not exist?”
It returns:
true → file does NOT exist
false → file exists, or some other error happened


basically when we run the final command

go run /Users/iwaili/api/cmd/main.go -config /Users/iwaili/api/config/local.yaml

this passes the configpath to main.go and it runs MustLoad, which finds this file and pass
the values of local.yaml to cfg struct of type Config
then this structs address' is returned
*/
