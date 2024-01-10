# JPC16 Core

## Command
- Clear all data
  ```
  go run -buildvcs=true ./_command clear
  ```
- Import player and group
  ```
  go run -buildvcs=true ./_command player import --file "C:\Users\BSthun\Downloads\Book1.csv"
  ```
- Export to CSV
  ```
  go run -buildvcs=true ./_command player export
  ```
- Start Discord Bot
  ```
  go run -buildvcs=true ./_bot
  ```