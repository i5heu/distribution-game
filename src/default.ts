export type Store = {
  golang: string;
};

export const defaults: Store = { "golang":"package main\n\nimport (\n    \"custom/custom\"\n\t\"strconv\"\n\t\"fmt\"\n    \"time\"\n)\n\nvar count = 0\n\nfunc test(){\n    for {\n        count++\n    \n time.Sleep(1 * time.Second)\n    }\n}\n\nfunc Receive(data custom.Data){\n    count++\n    custom.Send(data)\t\n}"}