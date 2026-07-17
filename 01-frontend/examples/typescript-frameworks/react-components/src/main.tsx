import { StrictMode, useState } from "react";
import { createRoot } from "react-dom/client";
import "./style.css";

interface Lesson {
  id: string;
  title: string;
  completed: boolean;
}

const lessons: readonly Lesson[] = [
  { id: "ts-01", title: "TypeScript 基础类型", completed: true },
  { id: "react-01", title: "React 组件与状态", completed: false },
  { id: "test-01", title: "前端测试层级", completed: false },
];

function LessonExplorer() {
  const [query, setQuery] = useState("");
  const [showCompleted, setShowCompleted] = useState(true);
  const normalized = query.trim().toLocaleLowerCase("zh-CN");
  const visible = lessons.filter((lesson) =>
    (showCompleted || !lesson.completed)
    && lesson.title.toLocaleLowerCase("zh-CN").includes(normalized)
  );

  return (
    <main>
      <section aria-labelledby="lesson-heading">
        <h1 id="lesson-heading">课程筛选</h1>
        <label>搜索<input value={query} onChange={(event) => setQuery(event.currentTarget.value)} /></label>
        <label className="checkbox"><input type="checkbox" checked={showCompleted} onChange={(event) => setShowCompleted(event.currentTarget.checked)} />显示已完成</label>
        <p aria-live="polite">{visible.length} 项</p>
        {visible.length === 0 ? <p>没有匹配课程</p> : (
          <ul>{visible.map((lesson) => <li key={lesson.id}>{lesson.title}<span>{lesson.completed ? "已完成" : "学习中"}</span></li>)}</ul>
        )}
      </section>
    </main>
  );
}

const container = document.getElementById("root");
if (!container) throw new Error("缺少 #root");
createRoot(container).render(<StrictMode><LessonExplorer /></StrictMode>);
