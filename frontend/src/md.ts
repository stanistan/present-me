import showdown from "showdown";

const converter = new showdown.Converter();
converter.setFlavor("github");

export const mdToHtml = (input: string): string => {
  return converter.makeHtml(input);
};
