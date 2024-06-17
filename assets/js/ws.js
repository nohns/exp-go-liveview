import morphdom from "./morphdom/morphdom.esm.js";

let ws;
let viewMap = new Map();
window.addEventListener("DOMContentLoaded", () => {
  const sid = document.body.dataset.gliveSession;
  ws = new WebSocket(`${window.location.origin}/ws/${sid}`);

  ws.onopen = function () {
    console.log("ws open!");
  };

  ws.onclose = function () {
    console.log("ws closed!");
  };

  ws.onerror = function (ev) {
    console.error("ws err!", ev);
  };

  ws.onmessage = function (ev) {
    const msg = JSON.parse(ev.data);
    const viewElement = document.querySelector(`[data-glive-view=${msg.vid}]`);
    const { s, ...values } = msg.data;
    const html = zip(s, values);

    morphdom(viewElement, html);
  };
});

function zip(segments, values) {
  if (segments.length === 1) {
    return segments[0];
  }

  let res = "";
  const valuesExpected = segments.length - 1;
  for (let i = 0; i < Object.keys(values).length; i++) {
    const segmentIndex = i % valuesExpected; // SegmentIndex of segment preceeding value
    if (segmentIndex == 0) {
      res += segments[0];
    }

    // Add value to result. Zip if value itself is a nested object
    let v = values[i];
    if (v) {
      if (typeof v === "object") {
        const { s: subSegments, ...subValues } = v;
        v = zip(subSegments, subValues);
      }
      res += v;
    }
    // Add value postfix
    res += segments[segmentIndex + 1];
  }
  return res;
}
