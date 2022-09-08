# fireblocks-go-sdk
Fireblocks GO SDK


## Description
This project is a library to use Fireblocks.

## Usage
Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

```golang
	fb, err := sdk.CreateSDK(
		apiKey,
		apiSecretKey,
		baseURL,
		sdk.WithTokenTimeout(5),
	)
```