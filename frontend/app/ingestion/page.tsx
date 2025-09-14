"use client";

import { useState } from "react";

export default function IngestionPage() {
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
      setMessage(null);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setMessage("Please select a CSV file first.");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      setUploading(true);
      setMessage(null);

      const res = await fetch("http://localhost:8080/ingest_csv", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        throw new Error(`Upload failed: ${res.statusText}`);
      }

      const json = await res.json();
      setMessage(`Upload successful: ${json.rows_ingested || "Unknown"} rows`);
      setIsSuccess(true);
      setFile(null);
    } catch (err: any) {
      setMessage(`Error: ${err.message}`);
      setIsSuccess(false);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-900 text-white">
      <div className="w-full max-w-md space-y-6 text-center">
        <h1 className="text-2xl font-bold">CSV Ingestion</h1>
        <p className="text-gray-400">
          Upload a CSV file to ingest into the database.
        </p>

        <div className="bg-gray-800 p-6 rounded-xl shadow space-y-4">
          <input
            type="file"
            accept=".csv"
            onChange={handleFileChange}
            className="block w-full text-sm text-gray-300
                       file:mr-4 file:py-2 file:px-4
                       file:rounded-full file:border-0
                       file:text-sm file:font-semibold
                       file:bg-blue-600 file:text-white
                       hover:file:bg-blue-700"
          />

          {file && (
            <p className="text-gray-300 text-sm">
              Selected file: <span className="font-semibold">{file.name}</span>
            </p>
          )}

          <button
            onClick={handleUpload}
            disabled={uploading}
            className={`px-4 py-2 rounded w-full ${
              uploading
                ? "bg-gray-500 cursor-not-allowed"
                : "bg-blue-600 hover:bg-blue-700"
            }`}
          >
            {uploading ? "Uploading..." : "Upload"}
          </button>

          {message && (
            <p
              className={`text-sm ${
                isSuccess ? "text-green-400" : "text-red-400"
              }`}
            >
              {message}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
