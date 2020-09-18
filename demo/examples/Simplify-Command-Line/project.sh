# Full sysl command line
sysl codegen --grammar=go.gen.g --transform=go.gen.sysl --start=goFile --app-name="Hello world"  model.sysl

# Can be simplified to 
sysl codegen @config.txt

# And config.txt has following command flags
--grammar=go.gen.g --transform=go.gen.sysl --start=goFile --app-name="Hello world" model.sysl