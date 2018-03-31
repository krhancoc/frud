# Frud - Framework for CRUD (Create, Read, Update, Delete) applications

In my attempt to learn Go I wanted to create a pluggable framework where the only work required by the user is to create a plugin like object and pass this to the server.  Another option is to create
their data models from a configuration object in JSON and pass this to the server, it will then create these models use the attached database drivers

Database's currently supported (in limited capacity):
* Neo4J
* MongoDB


This is by no means finished in any capacity, but will gladly accept pull requests.
## Getting Started
First thing is to compile the project
```bash
go build
```

### Plugin Definition Method
Next thing to do is to compile your plugins,  you can find an example plugin within the `_plugins` directory.  The reason for the underscore is that the go compiler will ignore directorys within the plugin.  You have to make sure to follow the rules dictated by the [go plugin library](https://golang.org/pkg/plugin/).  Make sure in the configuration object to set the `pathtocompiled`, and the `pathtocode`. The server does need both (look to the file "config.json" for an example config").

### Model Definition Method
To use the model method you are going to have to define it within the config.
```javascript
"plugins": [
      {
        "name": "model", //Required and unique
        "path": "/modelonly", //Required and unique
        "description": "Model only endpoint",
        "model": [
          {
            "key": "name", //Required and unique,
            "value_type": "string", // Required - values can be int, string,
            "options": ["id"] // Not required - must exist at most one model field with "id" option within the top layer of fields. Meaning ID cannot be within a nested field.  Although future option will be to allow a whole nested field to be an id.
          }
          {
            "key": "anotherExample", //Required and unique,
            "value_type": "int" //Required - values can be int, string,
            "options" : ["empty"] // Empty option allows for this field to be empty in a request - default is nonempty.
          }
          {
            "key": "deepValue", 
            "value_type" : [ // Can create nested types.
              {
                "key": "subfield_0", // Must be unique up to the nested field.
                "value_type": "another" // Can use user defined types
              },
              {
                "key": "subfield_1",
                "value_type": "string"
              }
            ]
          }
        ]
      },
      {
        "name": "another", //Required and unique
        "path": "/another", //Required and unique
        "description": "Model only endpoint",
        "model": [
          {
            "key": "name", //Required and unique,
            "value_type": "string", // Required - values can be int, string,
            "options": ["id"] // Not required - must exist at most one model field with "id" option
          }
        ]
      }
    ]
```
<!--
## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 
 -->
## Authors

* **Kenneth R Hancock** 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE) file for details

