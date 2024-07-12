build:
	@go build -o bin/app .

css:
	tailwindcss -i views/styles/main.css -o public/assets/styles/main.css --watch --minify
