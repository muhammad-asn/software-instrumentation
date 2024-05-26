### Nginx

LogQL: 
```bash
{app="nginx"} |= "GET / " 
| regexp `(?P<ip>\S+) (?P<identd>\S+) (?P<user>\S+) \[(?P<timestamp>[\w:\/]+\s[+\\-]\d{4})\] "(?P<action>\S+)\s?(?P<path>\S+)\s?(?P<protocol>\S+)?" (?P<status>\d{3}|-) (?P<size>\d+|-)\s?"?(?P<referrer>[^\"]*)"?\s?"?(?P<useragent>[^\"]*)?"?`
```