import React from 'react';

function MetalButton({ onClick, children }) {
  return (
    <button 
      style={{
        padding: '2rem 4rem',
        fontSize: '3rem',
        minWidth: '400px',
        borderRadius: '2rem',
        border: '4px solid #aaaaaa',
        fontWeight: 'bold',
        color: 'black',
        background: 'linear-gradient(145deg, #e6e6e6, #ffffff, #d4d4d4, #f5f5f5)',
        boxShadow: '0 8px 16px rgba(0,0,0,0.3), inset 0 1px 0 rgba(255,255,255,0.8), inset 0 -1px 0 rgba(0,0,0,0.2)',
        cursor: 'pointer',
        transition: 'transform 0.2s',
        fontFamily: "'Courier New', Courier, monospace"
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'scale(1.05)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'scale(1)';
      }}
      onClick={onClick}
    >
      {children}
    </button>
  );
}

export default MetalButton;
