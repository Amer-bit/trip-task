# Trip Task

### Development command
go mod tidy
go run .

### Build command
go build

### API's
- POST /trip
    - Body:
        - locationForm: coordinate | long-lat, "required", Example: "coordinate"
        - origin: string, "required", Example: "Husban & Um al Basateen, Amman"
        - destination:  string, "required", Example: "32.3949968, 35.9043093"

- GET /trip/price?origin=''&destination=''
- GET /trip/status/:status
    - status : pending | completed


# Deployed Domain
