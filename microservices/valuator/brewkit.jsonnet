local app = "service";

local copy = std.native('copy');

{
    apiVersion: "brewkit/v1",
    targets: {
        all: ['gobuild'],

        gobuild: {
            from: "golang:1.23",
            workdir: "/app",
            copy: [
                copy('cmd', 'cmd'),
                copy('pkg', 'pkg'),
            ],
            command: std.format("go build -o ./bin/%s ./cmd/%s", [app])
        }
    }
}