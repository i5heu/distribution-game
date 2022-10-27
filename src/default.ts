export type Store = {
  javascript: string;
};

export const defaults: Store = {
  javascript: "test('hello world') \n\n\n setInterval(function () {test('hello world')}, 1000);",
};
