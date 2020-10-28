# How to add your plugin as commands to odfe-cli

## Concepts

odfe-cli is built on a structure of 
1. Commands, arguments & flags
2. Controller
3. Entity
4. Gateway
5. Mapper
6. Handler

### 1. Commands, Arguments & Flags
Command is the central point of the application. Every interaction that cli 
supports will be represented as command. A command can have sub-commands, 
arguments and flags . Usually command represent an action, arguments are the 
things and flags are ways to modify the behavior of the commands.

A flag is a way to modify or tune the behavior of a command

In the following example, 'odfe-cli', 'ad', 'start' are commands, and 'profile' is a flag
```
odfe-cli ad start detector-name --profile prod
```

### 2. Controller
Controller implements actual business logic for a command. Every plugin should 
have at least one controller. A controller can depend on another controller and gateway. 

### 3. Entity
Model or Entities represents the request and response structure required to 
pass data to and from clusters and from commands to controller.

### 4. Gateway
Gateway acts like an interface between application and cluster. 
Every plugin will have a gateway which implements the REST api provided by 
corresponding plugin.

### 5. Mapper
This will perform operations like converting user input to REST API input format and REST API
output to user understandable format.

### 6. Handler
Since a controller is independent of commands, handler supports serialization and deserialization
of user input before it is passed on to controller.

## Getting started

odfe-cli follows the following organization structure:

```
▾ odfe-cli/
    ▾ commands/
        root.go
        root_test.go
        profile.go
        ad.go
        profile_test.go
    ▾ controller/
        ▾ ad/
            mocks/
            ad.go
            ad_test.go
    ▾ entity/
        ▾ ad/
            ad.go
            ad_test.go
    ▾ gateway/
        ▾ ad/
            mocks
            ad.go
            ad_test.go
    ▾ mapper/
        ▾ ad/
            testdata
            ad.go
            ad_test.go
    ▾ handler/
        ▾ ad/
            testdata
            ad.go
            ad_test.go
      main.go
```

## Integrate new plugins with odfe-cli
To integrate new plugins, you need to create a plugin base command file (odfe-cli/commands/ad.go). You will optionally provide additional commands as you see fit.
### Create plugin base command

```
//adCommand is base command for Anomaly Detection plugin.
var adCommand = &cobra.Command{
	Use:   "ad",
	Short: "Manage the Anomaly Detection plugin",
	Long: fmt.Sprintf("Description:\n  " +
		`Use the Anomaly Detection commands to create, configure, and manage detectors.`),
}
```
You can define flags and handle configuration in command's init() function.

```
func init() {
	adCommand.Flags().StringP(flagProfileName, "p", "", "Use a specific profile from your configuration file.")
	adCommand.Flags().BoolP("help", "h", false, "Help for Anomaly Detection")
	GetRoot().AddCommand(adCommand)
}
```

In above example, we created a base command 'ad' and added it to root command.

### Create additional commands

Additional commands can be defined and added in their own file inside commands/ directory.

If you wanted to implement a create detector command from Anomaly Detection, you would create odfe-cli/commands/ad_create.go
and include the following code snippet.

```
//createDetectorsCmd creates detectors
var createDetectorsCmd = &cobra.Command{
	Use:   createDetectorsCommandName + " json-file-path ...",
	Short: "Create detectors based on JSON files",
	Long: fmt.Sprintf("Description:\n  " +
		"Create detectors based on a local JSON file"),
	Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("create detector")
	},
}
func init() {
	createDetectorsCmd.Flags().StringP(flagProfileName, "p", "", "Use a specific profile from your configuration file.")
	createDetectorsCmd.Flags().BoolP("help", "h", false, "Help for "+createDetectorsCommandName)
	GetADCommand().AddCommand(createDetectorsCmd)
}
```
### Define Entity

Input and output for plugin's REST API is represented as entity (structure). For example, create detector REST API, needs
following json as input.

```
{
 "name": "test-detector",
 "description": "Test detector",
 "time_field": "timestamp",
 "indices": [
   "order*"
 ],
 "feature_attributes": [
   {
     "feature_name": "total_order",
     "feature_enabled": true,
     "aggregation_query": {
       "total_order": {
         "sum": {
           "field": "value"
         }
       }
     }
   }
 ],
 "filter_query": {},
 "detection_interval": {
   "period": {
     "interval": 1,
     "unit": "Minutes"
   }
 },
 "window_delay": {
   "period": {
     "interval": 1,
     "unit": "Minutes"
   }
 }
```
Dealing with JSON as it is error prone. Hence, convert it to structure like below and serialize or deserialize structures based
on use case using [json](https://golang.org/pkg/encoding/json) package.

```

//Feature structure for detector features
type Feature struct {
	Name             string          `json:"feature_name"`
	Enabled          bool            `json:"feature_enabled"`
	AggregationQuery json.RawMessage `json:"aggregation_query"`
}

//Period represents time interval
type Period struct {
	Duration int32  `json:"interval"`
	Unit     string `json:"unit"`
}

//Interval represent unit of time
type Interval struct {
	Period Period `json:"period"`
}

//CreateDetector represents Detector creation request
type CreateDetector struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	TimeField   string          `json:"time_field"`
	Index       []string        `json:"indices"`
	Features    []Feature       `json:"feature_attributes"`
	Filter      json.RawMessage `json:"filter_query,omitempty"`
	Interval    Interval        `json:"detection_interval"`
	Delay       Interval        `json:"window_delay"`
}

```
### Call Plugin REST API

Plugin's REST API can be abstracted as interface methods inside gateway. These methods will accept parameters if its
corresponding REST API expects payload, and returns output. 
```
// Gateway interface to AD Plugin
type Gateway interface {
	CreateDetector(context.Context, interface{}) ([]byte, error)
}
```
Gateway's only responsibility is to accept rest client in the constructor, build URL, create request, pass headers,
 and returns output.
```
func (g *gateway) CreateDetector(ctx context.Context, payload interface{}) ([]byte, error) {
	createURL, err := g.buildCreateURL()
	if err != nil {
		return nil, err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodPost, payload, createURL.String(), gw.GetHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(detectorRequest, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return response, nil
}
```

### Use Controller for business logic
Controller will act as connector between the handler and gateway. Every plugin will have a controller to
implement methods for every commands like below, for anomaly detection plugin.

Controller is the best place
to validate user input, transform user input for the gateway, map gateway response to user understandable format.

```
//Controller is an interface for the AD plugin controllers
type Controller interface {
	CreateAnomalyDetector(context.Context, string) error
}

type controller struct {
	gateway ad.Gateway
}

//CreateAnomalyDetector creates detector based on user request
func (c controller) CreateAnomalyDetector(ctx context.Context, r entity.CreateDetectorRequest) (*string, error) {

	if err := validateCreateRequest(r); err != nil {
		return nil, err
	}

        //convert user input to gateway input
	payload, err := mapper.MapToCreateDetector(r)
	if err != nil {
		return nil, err
	}

        // call gateway
	response, err := c.gateway.CreateDetector(ctx, payload)
	if err != nil {
		return nil, processEntityError(err)
	}

        // process gateway output
	var data map[string]interface{}
	_ = json.Unmarshal(response, &data)

	detectorID := fmt.Sprintf("%s", data["_id"])
        return mapper.StringToStringPtr(detectorID), nil
}

```

###  Use handler to bridge user commands and controller

Create a handler to call controller method accordingly. Most of the time, user's input and controller's method
will not be similar. For example: to create a detector, user will pass JSON file but controller accepts an entity, this is created
to keep controller logic agnostic to user file input. Hence, if ad supports YAML file in the future, controller
will not be updated.

```
//CreateAnomalyDetector creates detector based on file configurations
func (h *Handler) CreateAnomalyDetector(fileName string) error {
	if len(fileName) < 1 {
		return fmt.Errorf("file name cannot be empty")
	}

        **// read file, convert it into structure**
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file %s due to %v", fileName, err)
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("failed to close json:", err)
		}
	}()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var request entity.CreateDetectorRequest
	err = json.Unmarshal(byteValue, &request)
	if err != nil {
		return fmt.Errorf("file %s cannot be accepted due to %v", fileName, err)
	}
	ctx := context.Background()

        **// Call controller method**
	names, err := h.CreateAnomalyDetector(ctx, request)
	if err != nil {
		return err
	}

        **// convert controller output to user understandable format**
	if len(names) > 0 {
		fmt.Printf("Successfully created %d detector(s)", len(names))
		fmt.Println()
		return nil
	}
	return err
}

```
### Update command to call handler method

Finally, connect the handler to command's Run method to execute action from start to end.

```
//createDetectorsCmd creates detectors
var createDetectorsCmd = &cobra.Command{
	Use:   createDetectorsCommandName + " json-file-path ...",
	Short: "Create detectors based on JSON files",
	Long: fmt.Sprintf("Description:\n  " +
		"Create detectors based on a local JSON file"),
	Run: func(cmd *cobra.Command, args []string) {
           // fmt.Println("create detector")
            	commandHandler, err := GetADHandler()
            	for _, name := range fileNames {
            		err = handler.CreateAnomalyDetector(commandHandler, name)
            		if err != nil {
            			fmt.Println("failed to create detector due to ", err)
                        break;
            		}
            	}
	},
}
``` 

