# HOWTO

Run the command and follow the progress stream:
```shell
go run ./...
```

Output:
```shell
event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0,"stage":"queued"}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.01,"stage":"active"}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.01,"stage":"preload","stage_progress":0}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.11,"stage":"preload","stage_progress":0.5}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.21,"stage":"preload","stage_progress":1}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.21,"stage":"generate","stage_progress":0}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.32,"stage":"generate","stage_progress":0.2}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.53,"stage":"generate","stage_progress":0.6}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.55,"stage":"generate","stage_progress":0.64}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.57,"stage":"generate","stage_progress":0.68}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.74,"stage":"generate","stage_progress":1}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.74,"stage":"postprocessing","stage_progress":0}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.82,"stage":"postprocessing","stage_progress":0.33}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.91,"stage":"postprocessing","stage_progress":0.67}

event: generating
data: {"id":"4qiAHtediaz6i1tPbk","progress":0.99,"stage":"postprocessing","stage_progress":1}

event: completed
data: {"id":"xyz","progress":1,"stage":"complete","url":"https://peregrine-results.s3.amazonaws.com/pigeon/abcdefg.mp3","duration":1.2053,"size":25965}
```
