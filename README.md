# **gobuildproxy**

A build proxy that tracks changes and rebuilds if needed when a new request comes in. Build errors are captured and rendered in a more readable format right in your browser.
Screenshots are available below.

- 🚀 **Faster development workflow** – No need to wait for unnecessary builds.
- 🔋 **Saves CPU & battery** – Rebuilds only when required.
- 🔄 **Ensures latest changes are applied** – Refresh once and get the newest version.
- 🛠 **Build error reporting** – Go/Templ build errors directly in your browser.


## **🚀 Installation**
```sh
go install github.com/tigrang/gobuildproxy@latest
```

### **Available Flags**
| Flag          | Default Value    | Description                                             |
|---------------|------------------|---------------------------------------------------------|
| `--path`      | `.`              | Path to the app.                                        |
| `--run`       | `./run`          | Path to the app startup script.                         |
| `--cmd`       | `./build`        | Command to execute for rebuilding the application.      |
| `--proxybind` | `localhost:9000` | Address where gobuildproxy listens for requests.        |
| `--proxy`     | `localhost:3000` | URL of the application to forward requests to.          |
| `--timeout`   | `30`             | Time (seconds) to wait for the app to become available. |
| `--templ`     |                  | Enable templ proxy and error reporting.                 |

## **📝 Notes & Limitations**
- Your `run` script is responsible for stopping and starting your app in the background.

## **🛠 Usage**

### Go app (no Templ)
```sh
gobuildproxy
```

##### **📌 Example Configuration**

<details>
<summary>run script</summary>

```sh
#!/bin/sh

pkill myapp
./myapp &
```
</details>

<details>
<summary>build script</summary>
    
```sh
#!/bin/sh

go build -o myapp .
```
</details>

### Go app with Templ

gobuildproxy will start `templ generate --watch --proxy` to capture and render build errors.
Make sure to set `TEMPL_DEV_MODE=true` env var when starting your app.

```sh
gobuildproxy --templ
```

#### **📌 Example Configuration**

<details>
<summary>run script</summary>

```
#!/bin/sh

pkill myapp
TEMPL_DEV_MODE=true ./myapp &
```
</details>

<details>
<summary>build script</summary>

```
#!/bin/sh

go build -o myapp .
```
</details>

## **📷 Screenshots**
<img width="865" alt="Screenshot 2025-03-16 at 8 55 49 PM" src="https://github.com/user-attachments/assets/28699657-29e1-4d36-b7f2-9fd157f36abe" />

<img width="865" alt="Screenshot 2025-03-16 at 7 35 48 PM" src="https://github.com/user-attachments/assets/36f413f0-a3f3-431b-8861-7c7c28263a71" />

<img width="865" alt="Screenshot 2025-03-16 at 8 40 57 PM" src="https://github.com/user-attachments/assets/63fda4e5-735a-408a-94a0-dcece33cf17c" />


## **📜 License**
Licensed under the [MIT License](LICENSE).
