# easytime

## install
```
go install github.com/gamelife2025/easytime@latest
```
## Usage
```
easytime get 2024.1.29
# date: 2024-01-29T08:00:00+08:00
# timestamp: 1706486400
# weekday: 1
# zero_timestamp: 1706457600

easytime now
# date: 2024-01-29T13:54:00+08:00
# timestamp: 1706507640
# weekday: 1
# zero_timestamp: 1706457600
```
## web
```
easytime web --port 8000
```
![easytime web interface](./image.png)
## 参考实现
- github.com/arrow-py/arrow