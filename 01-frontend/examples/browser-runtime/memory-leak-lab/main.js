const sandbox = document.querySelector("#sandbox");
const metrics = {
  leakRuns: 0,
  safeRuns: 0,
  listeners: 0,
};

const leakedNodes = [];
const leakedListeners = [];

function createBatch() {
  const root = document.createElement("section");
  for (let index = 0; index < 200; index += 1) {
    const row = document.createElement("div");
    row.textContent = `row ${index}`;
    row.payload = new Array(500).fill(`payload-${index}`);
    root.append(row);
  }
  return root;
}

function renderMetrics(message) {
  document.querySelector("#leak-runs").textContent = metrics.leakRuns;
  document.querySelector("#safe-runs").textContent = metrics.safeRuns;
  document.querySelector("#retained").textContent = leakedNodes.length * 200;
  document.querySelector("#listeners").textContent = metrics.listeners;
  document.querySelector("#log").textContent =
    `${new Date().toLocaleTimeString()} ${message}\n` +
    document.querySelector("#log").textContent;
}

function runLeakingRound() {
  const root = createBatch();
  sandbox.append(root);
  const listener = () => root.dataset.lastResize = String(innerWidth);
  window.addEventListener("resize", listener);
  leakedListeners.push(listener);
  leakedNodes.push(root);
  root.remove();
  metrics.leakRuns += 1;
  metrics.listeners += 1;
  renderMetrics("泄漏轮次完成：全局数组和 window listener 仍保留 detached tree");
}

function runSafeRound() {
  const controller = new AbortController();
  const root = createBatch();
  sandbox.append(root);
  window.addEventListener(
    "resize",
    () => root.dataset.lastResize = String(innerWidth),
    { signal: controller.signal },
  );
  controller.abort();
  root.remove();
  metrics.safeRuns += 1;
  renderMetrics("安全轮次完成：监听已取消，局部引用可随函数返回释放");
}

function clearLeaks() {
  for (const listener of leakedListeners.splice(0)) {
    window.removeEventListener("resize", listener);
  }
  leakedNodes.length = 0;
  metrics.listeners = 0;
  renderMetrics("已删除实验中的全部长期引用；可在 Memory 面板请求 GC 后复测");
}

document.querySelector("#leak").addEventListener("click", runLeakingRound);
document.querySelector("#safe").addEventListener("click", runSafeRound);
document.querySelector("#clear").addEventListener("click", clearLeaks);
renderMetrics("实验就绪");
