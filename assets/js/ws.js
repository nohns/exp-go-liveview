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
    const viewElement = document.querySelector(
      `[data-glive-view=${msg.data.vid}]`,
    );
    switch (msg.type) {
      case "hydration":
        console.log(zip(msg.data.h));
        viewMap.set(msg.data.vid, msg.data.h);

        viewElement.outerHTML = zip(msg.data.h);
        break;
      case "diff":
        console.log("updated values", msg.data.v);
        viewElement.outerHTML = zip([viewMap.get(msg.data.vid)[0], msg.data.v]);
        break;
    }
  };
});

function zip(data) {
  let res = "";
  for (let i = 0; i < data[0].length * 2 - 1; i++) {
    res += data[i % 2][parseInt(i / 2)];
  }
  return res;
}
