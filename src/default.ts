export type Store = {
  golang: string;
};

export const defaults: Store = {
  "golang":"package main\n\nimport (\n    \"dg/dg\"\n    \"time\"\n    \"fmt\"\n)\n\nvar Info dg.Info\nvar terminate bool = false\n\nfunc Receive(data dg.Data){\n    //fmt.Println(\"DATA\", data, \"ME\", Info)\n}\n\nfunc Init(){\n    Info = dg.GetInfo()\n\n    for {\n        if terminate {\n            \n            break\n        }\n\t\t\ttime.Sleep(time.Second)\n    dg.Send(Info.Seeds[0], \"helllo\")\n\t\t}\n\n}\n\nfunc Close(){\n    terminate = true\n}"
};