"use client";

import { useState, useEffect, useRef } from "react";

type Job = {
  job_id: string;
  type: string;
  payload: string;
  status: string;
  created_at: string;
  updated_at: string;
};

const statusConfig: Record<string, { label: string; color: string; bg: string; bar: string }> = {
  pending:   { label: "待機中", color: "text-yellow-700", bg: "bg-yellow-50 border-yellow-200", bar: "bg-yellow-300" },
  running:   { label: "処理中", color: "text-blue-700",   bg: "bg-blue-50 border-blue-200",     bar: "bg-blue-400"   },
  completed: { label: "完了",   color: "text-green-700",  bg: "bg-green-50 border-green-200",   bar: "bg-green-400"  },
  failed:    { label: "失敗",   color: "text-red-700",    bg: "bg-red-50 border-red-200",       bar: "bg-red-400"    },
};

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

export default function Home() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitCount, setSubmitCount] = useState(6);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  // 今回投入したJob IDを記録するref
  const currentJobIdsRef = useRef<string[]>([]);

  // 今回投入したJob IDだけ個別に取得
  const fetchCurrentJobs = async () => {
    const ids = currentJobIdsRef.current;
    if (ids.length === 0) return;

    const results = await Promise.all(
      ids.map((id) =>
        fetch(`${apiUrl}/jobs/${id}`).then((r) => r.json())
      )
    );
    setJobs(results);
  };

  // 0.5秒ごとにポーリング
  const startPolling = () => {
    if (intervalRef.current) return;
    intervalRef.current = setInterval(fetchCurrentJobs, 500);
  };

  const stopPolling = () => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  };

  // 全Jobが完了・失敗したらポーリング停止
  useEffect(() => {
    if (jobs.length === 0) return;
    const allDone = jobs.every(
      (j) => j.status === "completed" || j.status === "failed"
    );
    if (allDone) stopPolling();
  }, [jobs]);

  useEffect(() => {
    return () => stopPolling();
  }, []);

  // Job一括投入
  const handleBulkSubmit = async () => {
    setIsSubmitting(true);
    setJobs([]);
    currentJobIdsRef.current = []; // IDをリセット
    stopPolling();

    try {
      const promises = Array.from({ length: submitCount }, (_, i) =>
        fetch(`${apiUrl}/jobs`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            type: "test",
            payload: `job-payload-${i + 1}`,
          }),
        })
      );
      const responses = await Promise.all(promises);
      const results = await Promise.all(responses.map((r) => r.json()));

      // 今回投入したJob IDだけ記録
      currentJobIdsRef.current = results.map((r) => r.job_id);

      // ポーリング開始
      startPolling();
    } catch (e) {
      console.error(e);
    } finally {
      setIsSubmitting(false);
    }
  };

  // ステータスごとの集計
  const counts = {
    pending:   jobs.filter((j) => j.status === "pending").length,
    running:   jobs.filter((j) => j.status === "running").length,
    completed: jobs.filter((j) => j.status === "completed").length,
    failed:    jobs.filter((j) => j.status === "failed").length,
  };

  const allDone =
    jobs.length > 0 &&
    jobs.every((j) => j.status === "completed" || j.status === "failed");

  return (
    <main className="min-h-screen bg-slate-50">
      {/* ヘッダー */}
      <header className="bg-white border-b border-slate-200 px-8 py-5 shadow-sm">
        <h1 className="text-2xl font-bold text-slate-800">
          Job Queue Dashboard
        </h1>
        <p className="text-sm text-slate-500 mt-1">
          並行処理のリアルタイム可視化
        </p>
      </header>

      <div className="max-w-4xl mx-auto px-6 py-10 space-y-8">
        {/* コントロールパネル */}
        <section className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
          <h2 className="text-lg font-semibold text-slate-700 mb-2">
            Job を一括投入する
          </h2>
          <p className="text-sm text-slate-500 mb-6">
            複数のJobを同時に投入して、Worker 3つが並行処理する様子を確認できます
          </p>

          <div className="flex items-center gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">
                投入数
              </label>
              <select
                value={submitCount}
                onChange={(e) => setSubmitCount(Number(e.target.value))}
                className="border border-slate-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
              >
                {[3, 6, 9, 12].map((n) => (
                  <option key={n} value={n}>{n}個</option>
                ))}
              </select>
            </div>

            <button
              onClick={handleBulkSubmit}
              disabled={isSubmitting}
              className="mt-5 bg-blue-500 hover:bg-blue-600 disabled:bg-slate-300 disabled:cursor-not-allowed text-white font-medium px-8 py-2 rounded-lg transition-colors text-sm"
            >
              {isSubmitting ? "投入中..." : "一括投入する"}
            </button>
          </div>
        </section>

        {/* ステータス集計 */}
        {jobs.length > 0 && (
          <div className="grid grid-cols-4 gap-4">
            {(["pending", "running", "completed", "failed"] as const).map((s) => {
              const cfg = statusConfig[s];
              return (
                <div
                  key={s}
                  className={`bg-white rounded-xl border p-4 text-center ${cfg.bg}`}
                >
                  <p className={`text-2xl font-bold ${cfg.color}`}>
                    {counts[s]}
                  </p>
                  <p className={`text-xs font-medium mt-1 ${cfg.color}`}>
                    {cfg.label}
                  </p>
                </div>
              );
            })}
          </div>
        )}

        {/* Job処理状況 */}
        {jobs.length > 0 && (
          <section className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold text-slate-700">
                Job 処理状況
              </h2>
              {allDone ? (
                <span className="text-sm text-green-600 font-medium bg-green-50 px-3 py-1 rounded-full border border-green-200">
                  ✅ 全Job完了
                </span>
              ) : (
                <span className="text-sm text-blue-600 font-medium bg-blue-50 px-3 py-1 rounded-full border border-blue-200 animate-pulse">
                  ⚡ 処理中...
                </span>
              )}
            </div>

            <div className="space-y-3">
              {jobs.map((job) => {
                const cfg = statusConfig[job.status] ?? statusConfig.pending;
                return (
                  <div
                    key={job.job_id}
                    className={`border rounded-xl p-4 transition-all duration-300 ${cfg.bg}`}
                  >
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-xs font-mono text-slate-400 truncate w-72">
                        {job.job_id}
                      </span>
                      <span className={`text-xs font-semibold px-2 py-0.5 rounded-full border ${cfg.bg} ${cfg.color}`}>
                        {cfg.label}
                      </span>
                    </div>

                    {/* ステータスバー */}
                    <div className="w-full bg-slate-100 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full transition-all duration-500 ${cfg.bar} ${
                          job.status === "running" ? "animate-pulse" : ""
                        }`}
                        style={{
                          width:
                            job.status === "pending"   ? "5%"  :
                            job.status === "running"   ? "60%" : "100%",
                        }}
                      />
                    </div>

                    <div className="flex justify-between mt-2">
                      <span className="text-xs text-slate-400">
                        payload: {job.payload}
                      </span>
                      <span className="text-xs text-slate-400">
                        更新: {job.updated_at}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          </section>
        )}
      </div>
    </main>
  );
}