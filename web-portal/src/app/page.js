/* eslint-disable @next/next/no-img-element */
"use client";

import React, { useState, useRef, useCallback, useEffect } from 'react';

// The main App component that handles the entire drag-and-drop functionality
export default function App() {
  // State to manage the drag-and-drop area's visual feedback
  const [isDragging, setIsDragging] = useState(false);
  // State to store the selected file
  const [file, setFile] = useState(null);
  // State to track the overall job status
  const [jobStatus, setJobStatus] = useState('idle');
  // State to track the render progress percentage
  const [progress, setProgress] = useState(0);
  // State for displaying a specific status message to the user
  const [statusMessage, setStatusMessage] = useState(''); // New state for custom messages
  // State for displaying error messages
  const [errorMessage, setErrorMessage] = useState('');
  // State for adding local ip of the server
  const [serverIp, setServerIp] = useState("http://192.168.x.x:5000");
  // State to manage the dropzone border color based on job status
  const [dropzoneBorderClass, setDropzoneBorderClass] = useState('border-zinc-800');
  // New state to hold the URL of the rendered image
  const [downloadUrl, setDownloadUrl] = useState('');


  // Reference to the hidden file input element
  const fileInputRef = useRef(null);

  // Effect to update dropzone border class when jobStatus changes
  useEffect(() => {
    if (jobStatus === 'success') {
      setDropzoneBorderClass('border-green-500');
    } else if (jobStatus === 'error') {
      setDropzoneBorderClass('border-red-500');
    } else {
      setDropzoneBorderClass('border-zinc-800');
    }
  }, [jobStatus]);

  // Simple SVG for the cloud icon to avoid external libraries
  const CloudUploadIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="60" height="60" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1" strokeLinecap="round" strokeLinejoin="round" className="text-gray-400 mb-4 drop-shadow-lg">
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
      setJobStatus('error');
      setFile(null);
      return;
    }

    // Reset state and set the file
    setErrorMessage('');
    setFile(uploadedFile);
    setJobStatus('idle');
    setProgress(0);
  };

  // The main function to handle file upload and rendering
  const uploadFile = async () => {
    if (!file) {
      setErrorMessage("No file selected to upload.");
      setJobStatus("error");
      return;
    }

    if (!serverIp) {
      setErrorMessage("Please provide your server's IP address.");
      setJobStatus("error");
      return;
    }

    try {
      setJobStatus("uploading");
      setStatusMessage("Uploading file...");
      setProgress(0); // Reset progress

      const formData = new FormData();
      formData.append("file", file);

      // Fetch API with upload progress
      const xhr = new XMLHttpRequest();
      xhr.open('POST', `${serverIp}/upload`, true);

      xhr.upload.onprogress = (e) => {
        if (e.lengthComputable) {
          const percent = Math.round((e.loaded / e.total) * 100);
          setProgress(percent);
          setStatusMessage(`Uploading: ${percent}%`);
        }
      };

      xhr.onloadend = async () => {
        if (xhr.status === 200) {
          // Upload is complete, now start polling for render status
          setJobStatus("rendering");
          setStatusMessage("File uploaded. Rendering started...");
          setProgress(0);

          const renderFileName = file.name.replace('.blend', '_render0001.png');

          while (true) {
            await new Promise(resolve => setTimeout(resolve, 1000)); // Poll every 1 seconds

            const statusResponse = await fetch(`${serverIp}/status?file=${file.name}`);
            const statusText = await statusResponse.text();

            console.log('Server status:', statusText);

            // shows a progress bar for rendering.
            if (statusText.includes("in-progress")) {
              setProgress(currentProgress => Math.min(currentProgress + 10, 99));
              setStatusMessage("Rendering in progress...");
            } else if (statusText.includes("done")) {
              setProgress(100);
              setJobStatus("success");
              setStatusMessage("Rendering completed!");
              // Construct the download URL and save it to state
              setDownloadUrl(`${serverIp}/download?file=${renderFileName}`);
              break; // Exit the loop
            } else if (statusText.includes("failed")) {
              setErrorMessage("Rendering failed on the server.");
              setJobStatus("error");
              break; // Exit the loop
            } else {
              setErrorMessage("An unexpected error occurred on the server.");
              setJobStatus("error");
              break;
            }
          }
        } else {
          setErrorMessage(xhr.responseText || "Upload failed.");
          setJobStatus("error");
        }
      };

      xhr.onerror = () => {
        setErrorMessage("Upload failed due to a network error.");
        setJobStatus("error");
      };

      xhr.send(formData);

    } catch (err) {
      setErrorMessage(err.message || "Upload failed.");
      setJobStatus("error");
    }
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
    setDropzoneBorderClass('border-amber-300'); // Set hover border during drag
  }, [handleDragEvents]);

  // Handler for when a dragged file leaves the drop zone
  const handleDragLeave = useCallback((e) => {
    handleDragEvents(e);
    setIsDragging(false);
    // Reset to default border if no specific status
    if (jobStatus === 'idle') {
      setDropzoneBorderClass('border-zinc-800');
    }
  }, [handleDragEvents, jobStatus]);

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
    switch (jobStatus) {
      case 'idle':
        return (
          <>
            <CloudUploadIcon />
            <p className="text-xl font-semibold mb-2 text-gray-100">Drag and drop a .blend file here</p>
            <p className="text-sm text-gray-400 mb-4">or</p>
            <button
              onClick={() => fileInputRef.current.click()}
              className="bg-zinc-700 text-white font-bold py-2 px-6 rounded-full shadow-lg hover:shadow-xl transition-all border-2 border-zinc-600 hover:bg-zinc-800"
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
      case 'rendering':
        return (
          <div className="w-full text-center">
            <p className="text-lg font-bold text-gray-100 mb-4">{statusMessage}</p>
            <div className="w-full bg-zinc-800 rounded-full h-2.5 overflow-hidden">
              <div
                className="bg-amber-300 h-2.5 rounded-full transition-all duration-300 ease-in-out shadow-inner"
                style={{ width: `${progress}%` }}
              ></div>
            </div>
            <p className="text-xs text-gray-400 mt-2">{progress}% Complete</p>
          </div>
        );
      case 'success':
        return (
          <div className="text-center text-green-500">
            <p className="text-4xl mb-4">✅</p>
            <p className="text-lg font-bold">Rendering Completed!</p>
            <p className="text-sm text-gray-400">Your file &quot;{file?.name}&quot; has been rendered successfully.</p>
            {downloadUrl && (
              <a 
                href={downloadUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="bg-green-700 text-white mt-4 font-bold py-2 px-6 rounded-full shadow-lg hover:bg-green-600 transition-colors border-2 border-green-600 inline-block mr-2"
              >
                View Render
              </a>
            )}
            <button
              onClick={() => {
                setJobStatus('idle');
                setDownloadUrl(''); // Clear the URL when resetting
              }}
              className="bg-zinc-700 text-white mt-4 font-bold py-2 px-6 rounded-full shadow-lg hover:bg-zinc-800 transition-colors border-2 border-zinc-600 inline-block"
            >
              Upload Another
            </button>
          </div>
        );
      case 'error':
        return (
          <div className="text-center text-red-400">
            <p className="text-4xl mb-4">❌</p>
            <p className="text-lg font-bold">Job Failed</p>
            <p className="text-sm text-gray-400">{errorMessage}</p>
            <button
              onClick={() => setJobStatus('idle')}
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
    <div className="font-sans flex flex-col items-center justify-center min-h-screen bg-zinc-950 text-gray-100 p-8 pb-20 gap-16 sm:p-20">
      <main className="max-w-xl w-full p-8 bg-zinc-900 rounded-3xl shadow-2xl backdrop-blur-sm bg-opacity-70 border border-zinc-800 transition-all duration-300 row-start-2 items-center sm:items-start">
        <h1 className="text-4xl font-extrabold text-center mb-6 text-transparent bg-clip-text bg-gradient-to-r from-gray-400 to-gray-200 drop-shadow-xl">
          Distributed Rendering
        </h1>

        {/* Server IP Input */}
        <div className="w-full mb-6">
          <h2 className="text-lg font-semibold text-gray-100 mb-2">Server IP</h2>
          <input
            type="text"
            value={serverIp}
            onChange={(e) => setServerIp(e.target.value)}
            placeholder="192.168.x.x:5000"
            className="w-full p-3 rounded-lg bg-zinc-800 text-amber-300 border border-zinc-700 focus:outline-none focus:border-amber-300 transition-colors shadow-inner"
          />
        </div>

        <div
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current?.click()}
          className={`
            flex flex-col items-center justify-center p-8 text-center
            border-4 border-dashed rounded-2xl cursor-pointer transition-colors duration-200
            ${isDragging
              ? 'border-amber-300 bg-zinc-800 bg-opacity-70'
              : dropzoneBorderClass + ' hover:border-amber-300 hover:bg-zinc-800'
            }
            shadow-lg
          `}
        >
          {renderContent()}
        </div>

        {/* Display file name and upload button if a file is selected */}
        {file && jobStatus === 'idle' && (
          <div className="mt-6 p-4 bg-zinc-800 rounded-lg flex flex-col sm:flex-row items-center justify-between shadow-lg border border-zinc-700">
            <span className="text-gray-300 font-medium truncate mb-2 sm:mb-0 sm:mr-4">
              File: {file.name}
            </span>
            <button
              onClick={uploadFile}
              className="bg-zinc-700 text-white font-bold py-2 px-6 rounded-full shadow-lg hover:bg-zinc-800 transition-colors"
            >
              Upload
            </button>
          </div>
        )}
      </main>
      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4 text-gray-300 transition-colors hover:text-amber-300"
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
