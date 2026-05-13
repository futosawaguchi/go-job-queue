"use client";

import { useState } from "react";

type Job = {
  job_id: string;
  type: string;
  payload: string;
  status: string;
  created_at: string;
  updated_at: string;
};

const statusColor: Record<string, string> = {
  pending: "bg-yellow-100 text-yellow-700",
  running: "bg-blue-100 text-blue-700",
  completed: "bg-green-100 text-green-700",
  failed: "bg-red-100 text-red-700",
};

export default function Home() {
  // Job投入フォームの状態
  const [type, setType] = useState("");
  const [payload, setPayload] = useState("");
  const [submitResult, setSubmitResult] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  // Job確認フォームの状態
  const [jobId, setJobId] = useState("");
  const [job, setJob] = useState<Job | null>(null);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  // POST /jobs
  const handleSubmit = async () => {
    setSubmitResult(null);
    setSubmitError(null);
    try {
      const res = await fetch(`${apiUrl}/jobs`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type, payload }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error);
      setSubmitResult(data.job_id);
      setType("");
      setPayload("");
    } catch (e: unknown) {
      setSubmitError(e instanceof Error ? e.message : "エラーが発生しました");
    }
  };

  // GET /jobs/:id
  const handleFetch = async () => {
    setJob(null);
    setFetchError(null);
    try {
      const res = await fetch(`${apiUrl}/jobs/${jobId}`);
      const data = await res.json();
      if (!res.ok) throw new Error(data.error);
      setJob(data);
    } catch (e: unknown) {
      setFetchError(e instanceof Error ? e.message : "エラーが発生しました");
    }
  };

  return (
    <main className="min-h-screen bg-slate-50">
      {/* ヘッダー */}
      <header className="bg-white border-b border-slate-200 px-8 py-5 shadow-sm">
        <h1 className="text-2xl font-bold text-slate-800">
          Job Queue Dashboard
        </h1>
        <p className="text-sm text-slate-500 mt-1">
          Jobの投入・ステータス確認ができます
        </p>
      </header>

      <div className="max-w-3xl mx-auto px-6 py-10 space-y-8">
        {/* Job投入カード */}
        <section className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
          <h2 className="text-lg font-semibold text-slate-700 mb-6">
            Job を投入する
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">
                Type
              </label>
              <input
                type="text"
                value={type}
                onChange={(e) => setType(e.target.value)}
                placeholder="例: image-resize"
                className="w-full border border-slate-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">
                Payload
              </label>
              <textarea
                value={payload}
                onChange={(e) => setPayload(e.target.value)}
                placeholder="例: {&quot;url&quot;: &quot;https://example.com/image.jpg&quot;}"
                rows={3}
                className="w-full border border-slate-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400 resize-none"
              />
            </div>
            <button
              onClick={handleSubmit}
              disabled={!type || !payload}
              className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-slate-300 disabled:cursor-not-allowed text-white font-medium py-2.5 rounded-lg transition-colors text-sm"
            >
              Job を投入する
            </button>

            {/* 成功メッセージ */}
            {submitResult && (
              <div className="bg-green-50 border border-green-200 rounded-lg px-4 py-3">
                <p className="text-sm text-green-700 font-medium">
                  ✅ Job を投入しました
                </p>
                <p className="text-xs text-green-600 mt-1 font-mono">
                  ID: {submitResult}
                </p>
              </div>
            )}

            {/* エラーメッセージ */}
            {submitError && (
              <div className="bg-red-50 border border-red-200 rounded-lg px-4 py-3">
                <p className="text-sm text-red-700">❌ {submitError}</p>
              </div>
            )}
          </div>
        </section>

        {/* Job確認カード */}
        <section className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
          <h2 className="text-lg font-semibold text-slate-700 mb-6">
            Job のステータスを確認する
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">
                Job ID
              </label>
              <input
                type="text"
                value={jobId}
                onChange={(e) => setJobId(e.target.value)}
                placeholder="例: fa751d62-4584-4bfe-89da-c4fec714c94b"
                className="w-full border border-slate-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400 font-mono"
              />
            </div>
            <button
              onClick={handleFetch}
              disabled={!jobId}
              className="w-full bg-slate-700 hover:bg-slate-800 disabled:bg-slate-300 disabled:cursor-not-allowed text-white font-medium py-2.5 rounded-lg transition-colors text-sm"
            >
              ステータスを確認する
            </button>

            {/* Job情報 */}
            {job && (
              <div className="bg-slate-50 border border-slate-200 rounded-lg p-4 space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-600">
                    ステータス
                  </span>
                  <span
                    className={`text-xs font-semibold px-3 py-1 rounded-full ${statusColor[job.status] ?? "bg-slate-100 text-slate-600"}`}
                  >
                    {job.status}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-600">
                    Type
                  </span>
                  <span className="text-sm text-slate-800">{job.type}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-600">
                    Payload
                  </span>
                  <span className="text-sm text-slate-800 font-mono">
                    {job.payload}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-600">
                    作成日時
                  </span>
                  <span className="text-sm text-slate-500">{job.created_at}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-600">
                    更新日時
                  </span>
                  <span className="text-sm text-slate-500">{job.updated_at}</span>
                </div>
                <div>
                  <span className="text-sm font-medium text-slate-600">
                    Job ID
                  </span>
                  <p className="text-xs text-slate-400 font-mono mt-1 break-all">
                    {job.job_id}
                  </p>
                </div>
              </div>
            )}

            {/* エラーメッセージ */}
            {fetchError && (
              <div className="bg-red-50 border border-red-200 rounded-lg px-4 py-3">
                <p className="text-sm text-red-700">❌ {fetchError}</p>
              </div>
            )}
          </div>
        </section>
      </div>
    </main>
  );
}