import _bytes from "bytes";
import * as timeago from "timeago.js";

timeago.register("en_short_no_ago", function (_: any, index: any): any {
  return [
    ["just now", "soon"],
    ["%ss", "%ss"],
    ["1m", "1m"],
    ["%sm", "%sm"],
    ["1h", "1h"],
    ["%sh", "%sh"],
    ["1d", "1d"],
    ["%sd", "%sd"],
    ["1w", "1w"],
    ["%sw", "%sw"],
    ["1mo", "1mo"],
    ["%smo", "%smo"],
    ["1yr", "1yr"],
    ["%syr", "%syr"],
  ][index];
});

export function prettyDate(date: Date): string {
  return timeago.format(date, "en_short_no_ago");
}

export function prettySize(bytes: number): string {
  return _bytes(bytes, {
    unitSeparator: " ",
    decimalPlaces: 1,
    fixedDecimals: true,
  });
}
