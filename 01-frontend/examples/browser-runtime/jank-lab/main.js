const status = document.querySelector("#status");
const bar = document.querySelector("#progress span");
const boxes = document.querySelector("#boxes");
const values = Array.from({ length: 1_500_000 }, (_, index) => index);

for (let index = 0; index < 400; index += 1) {
  const box = document.createElement("div");
  box.className = "box";
  boxes.append(box);
}

function expensive(value) {
  let result = value;
  for (let step = 0; step < 12; step += 1) {
    result = (result * 1.000001 + step) % 100003;
  }
  return result;
}

function syncWork() {
  const started = performance.now();
  let checksum = 0;
  for (const value of values) checksum += expensive(value);
  status.value = `同步完成 ${Math.round(performance.now() - started)} ms，${Math.round(checksum)}`;
}

function taskYield() {
  return new Promise((resolve) => {
    const channel = new MessageChannel();
    channel.port1.onmessage = () => resolve();
    channel.port2.postMessage(undefined);
  });
}

async function yieldedWork() {
  const started = performance.now();
  let checksum = 0;
  let index = 0;
  while (index < values.length) {
    const deadline = performance.now() + 8;
    while (index < values.length && performance.now() < deadline) {
      checksum += expensive(values[index]);
      index += 1;
    }
    bar.style.width = `${index / values.length * 100}%`;
    if (index < values.length) await taskYield();
  }
  status.value = `切片完成 ${Math.round(performance.now() - started)} ms，${Math.round(checksum)}`;
}

function thrashLayout() {
  const started = performance.now();
  for (const [index, box] of [...boxes.children].entries()) {
    box.style.width = `${box.offsetWidth + index % 4}px`;
  }
  status.value = `交错读写 ${Math.round(performance.now() - started)} ms`;
}

function batchLayout() {
  const started = performance.now();
  const widths = [...boxes.children].map((box) => box.offsetWidth);
  [...boxes.children].forEach((box, index) => {
    box.style.width = `${widths[index] + index % 4}px`;
  });
  status.value = `批量读写 ${Math.round(performance.now() - started)} ms`;
}

document.querySelector("#sync").addEventListener("click", syncWork);
document.querySelector("#yield").addEventListener("click", yieldedWork);
document.querySelector("#thrash").addEventListener("click", thrashLayout);
document.querySelector("#batch").addEventListener("click", batchLayout);
