import React from 'react';

function MetalContainer({ children, title }) {
  return (
    <div 
      style={{
        border: '4px solid #aaaaaa',
        borderRadius: '2rem',
        background: 'linear-gradient(145deg, #e6e6e6, #ffffff, #d4d4d4, #f5f5f5)',
        boxShadow: '0 8px 16px rgba(0,0,0,0.3), inset 0 1px 0 rgba(255,255,255,0.8), inset 0 -1px 0 rgba(0,0,0,0.2)',
        margin: '2rem auto',
        padding: '2rem',
        maxWidth: '800px'
      }}
    >
      {title && <h2 style={{
        fontSize: '1.5rem',
        fontWeight: 'bold',
        color: 'black',
        textAlign: 'center',
        marginBottom: '1rem'
      }}>{title}</h2>}
      <div style={{ margin: '1rem' }}>
        {children}
      </div>
    </div>
  );
}

export default MetalContainer;
