/* eslint-disable @next/next/no-img-element */
"use client";

import React, { useState, useRef, useCallback } from 'react';

// The main App component that handles the entire drag-and-drop functionality
export default function App() {
  // State to manage the drag-and-drop area's visual feedback
  const [isDragging, setIsDragging] = useState(false);
  // State to store the selected file
  const [file, setFile] = useState(null);
  // State to track the upload status: 'idle', 'uploading', 'success', 'error'
  const [uploadStatus, setUploadStatus] = useState('idle');
  // State for the upload progress percentage
  const [progress, setProgress] = useState(0);
  // State for displaying error messages
  const [errorMessage, setErrorMessage] = useState('');

  // Reference to the hidden file input element
  const fileInputRef = useRef(null);

  // Simple SVG for the cloud icon to avoid external libraries
  const CloudUploadIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="60" height="60" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1" strokeLinecap="round" strokeLinejoin="round" className="text-gray-400 mb-4">
      <path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242M12 12v9m-4-8 4-4 4 4"/>
    </svg>
  );

  // Handler for when a file is selected via the button
  const handleFileSelect = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      handleFile(selectedFile);
    }
  };

  // Main function to handle and validate the file
  const handleFile = (uploadedFile) => {
    // Check for file type
    if (!uploadedFile.name.endsWith('.blend')) {
      setErrorMessage('Invalid file type. Please upload a .blend file.');
      setFile(null);
      setUploadStatus('error');
      return;
    }

    // Reset state and set the file
    setErrorMessage('');
    setFile(uploadedFile);
    setUploadStatus('idle');
    setProgress(0);
  };

  // Simulates a file upload with a progress bar. Replace this with your actual API call.
  const uploadFile = () => {
    if (!file) {
      setErrorMessage('No file selected to upload.');
      setUploadStatus('error');
      return;
    }

    setUploadStatus('uploading');
    setProgress(0);

    // Mock progress over 3 seconds
    const interval = setInterval(() => {
      setProgress((prevProgress) => {
        if (prevProgress >= 100) {
          clearInterval(interval);
          setUploadStatus('success');
          return 100;
        }
        return prevProgress + 10;
      });
    }, 300);

  };

  // Event handler for drag-and-drop events (prevents default behavior)
  const handleDragEvents = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  // Handler for when a file is dragged over the drop zone
  const handleDragOver = useCallback((e) => {
    handleDragEvents(e);
    setIsDragging(true);
  }, [handleDragEvents]);

  // Handler for when a dragged file leaves the drop zone
  const handleDragLeave = useCallback((e) => {
    handleDragEvents(e);
    setIsDragging(false);
  }, [handleDragEvents]);

  // Handler for when a file is dropped
  const handleDrop = useCallback((e) => {
    handleDragEvents(e);
    setIsDragging(false);
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      handleFile(droppedFile);
    }
  }, [handleDragEvents]);

  // Renders the main UI based on the current state
  const renderContent = () => {
    switch (uploadStatus) {
      case 'idle':
        return (
          <>
            <CloudUploadIcon />
            <p className="text-xl font-semibold mb-2 text-gray-200">Drag and drop a .blend file here</p>
            <p className="text-sm text-gray-400 mb-4">or</p>
            <button
              onClick={() => fileInputRef.current.click()}
              className="bg-slate-700 text-white font-bold py-2 px-6 rounded-full shadow-lg hover:shadow-xl transition-shadow border-2 border-slate-600 hover:bg-slate-600"
            >
              Browse Files
            </button>
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleFileSelect}
              accept=".blend"
              className="hidden"
            />
          </>
        );
      case 'uploading':
        return (
          <div className="w-full text-center">
            <p className="text-lg font-bold text-gray-200 mb-4">Uploading: {file?.name}</p>
            <div className="w-full bg-gray-700 rounded-full h-2.5">
              <div
                className="bg-slate-500 h-2.5 rounded-full transition-all duration-300 ease-in-out"
                style={{ width: `${progress}%` }}
              ></div>
            </div>
            <p className="text-xs text-gray-400 mt-2">{progress}% Complete</p>
          </div>
        );
      case 'success':
        return (
          <div className="text-center text-green-400">
            <p className="text-4xl mb-4">✅</p>
            <p className="text-lg font-bold">Upload Successful!</p>
            <p className="text-sm text-gray-400">Your file &quot;{file?.name}&quot; is now rendering.</p>
            <button
              onClick={() => setUploadStatus('idle')}
              className="bg-slate-700 text-white mt-4 font-bold py-2 px-6 rounded-full shadow-lg hover:bg-slate-600 transition-colors border-2 border-slate-600"
            >
              Upload Another
            </button>
          </div>
        );
      case 'error':
        return (
          <div className="text-center text-red-400">
            <p className="text-4xl mb-4">❌</p>
            <p className="text-lg font-bold">Upload Failed</p>
            <p className="text-sm text-gray-400">{errorMessage}</p>
            <button
              onClick={() => setUploadStatus('idle')}
              className="bg-red-700 text-white mt-4 font-bold py-2 px-6 rounded-full shadow-lg hover:bg-red-600 transition-colors border-2 border-red-600"
            >
              Try Again
            </button>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="font-sans flex flex-col items-center justify-center min-h-screen bg-gray-950 text-gray-100 p-8 pb-20 gap-16 sm:p-20">
      <main className="max-w-xl w-full p-8 bg-gray-800 rounded-3xl shadow-2xl backdrop-blur-sm bg-opacity-70 border border-gray-700 transition-all duration-300 row-start-2 items-center sm:items-start">
        <h1 className="text-4xl font-extrabold text-center mb-6 text-transparent bg-clip-text bg-gradient-to-r from-slate-400 to-slate-600">
          Distributed Rendering
        </h1>

        {/* The main drop zone area */}
        <div
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current.click()}
          className={`
            flex flex-col items-center justify-center p-8 text-center
            border-4 border-dashed rounded-2xl cursor-pointer transition-colors duration-200
            ${isDragging
              ? 'border-slate-500 bg-slate-900 bg-opacity-30'
              : 'border-gray-600 hover:border-slate-500 hover:bg-gray-700'
            }
          `}
        >
          {renderContent()}
        </div>

        {/* Display file name and upload button if a file is selected */}
        {file && uploadStatus === 'idle' && (
          <div className="mt-6 p-4 bg-gray-700 rounded-lg flex flex-col sm:flex-row items-center justify-between">
            <span className="text-gray-300 font-medium truncate mb-2 sm:mb-0 sm:mr-4">
              File: {file.name}
            </span>
            <button
              onClick={uploadFile}
              className="bg-slate-500 text-white font-bold py-2 px-6 rounded-full shadow-lg hover:bg-slate-600 transition-colors"
            >
              Upload
            </button>
          </div>
        )}
      </main>
      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4"
          href="https://www.0shree005.tech/"
          target="_blank"
          rel="noopener noreferrer"
        >
          <img
            aria-hidden
            src="/earth.svg"
            alt="Globe icon"
            width={16}
            height={16}
          />
          About me: 0shree005.tech →
        </a>
      </footer>
    </div>
  );
}

