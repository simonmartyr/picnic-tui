# Picnic-Tui

An unofficial terminal interface for the online supermarket Picnic.

![demo](./screenshots/demo.gif)

## Features

- Search Picnic.
- Add/Remove articles to your basket.
- Browse and pick your delivery time.
- Confirm and pay for order (via ideal link)

## Running the app

### Windows

Download the 32 or 64 bit exe and via the cmd/powershell run:

`./picnic-tui-win64.exe -t <your auth token>`

alternately

`./picnic-tui-win64.exe -u <your username> -p <your password>`

`./picnic-tui-win64.exe -u <your username> -hp <your md5 hashed password>`

### Linux / Mac

Download the `picnic_tui_x32` or `picnic_tui_x64` for linux or `picnic_tui_mac` for mac and run via the commandline:

`./picnic_tui_x64 -t <your auth token>`

alternately

`./picnic_tui_x64 -u <your username> -p <your password>`

`./picnic_tui_x64 -u <your username> -hp <your md5 hashed password>`


## Keybindings 

The demo above highlights most but here is a complete breakdown:

| Location              | Operation                         | Binding                        |
|-----------------------|-----------------------------------|--------------------------------|
| Global                | Vim Style Exit                    | <kbd>:</kbd> <kbd>q</kbd>      |
|                       |                                   |                                |
| Main page             | Refresh                           | <kbd>Ctrl</kbd> + <kbd>R</kbd> |
| Main page             | Search (clear text)               | <kbd>Ctrl</kbd> + <kbd>S</kbd> |
| Main Page             | Search                            | <kbd>S</kbd> or <kbd>/</kbd>   |
| Main Page             | See more information on a product | <kbd>f</kbd>                   |
| Main Page             | Add 1 Item                        | <kbd>Enter</kbd>               |
| Main Page             | Add x Items                       | <kbd>0</kbd> - <kbd>9</kbd>    |
| Main Page             | Remove 1 Item                     | <kbd>backspace</kbd>           |
| Main Page             | Clear Basket                      | <kbd>Ctrl</kbd> + <kbd>K</kbd> |
| Main Page             | Switch to Delivery Page           | <kbd>D</kbd>                   |
| Main Page             | Switch to Checkout Page           | <kbd>C</kbd>                   |
| Main Page             | Switch to Delivery Tracker Page   | <kbd>T</kbd>                   |
|                       |                                   |                                |
| Delivery Page         | Return to Main Page               | <kbd>Esc</kbd>                 |
| Delivery Page         | Select Delivery Slot              | <kbd>Enter</kbd>               |
|                       |                                   |                                |
| Checkout Page         | Return to Main Page               | <kbd>Esc</kbd>                 |
| Checkout Page         | Start checkout process            | <kbd>c</kbd>                   |
|                       |                                   |                                |
| Delivery Tracker Page | Return to Main Page               | <kbd>Esc</kbd>                 |
| Delivery Tracker Page | Select Order                      | <kbd>Enter</kbd>               |
| Delivery Tracker Page | Navigate                          | <kbd>Tab</kbd>                 |
