import React, { useState } from 'react';

function LinkButton({ onClick, children, disabled = false }) {
  const [isHovered, setIsHovered] = useState(false);

  const buttonStyle = {
    width: isHovered && !disabled ? '99%' : '100%',
    height: isHovered && !disabled ? '98%' : '100%',
    background: disabled
      ? 'linear-gradient(145deg, #999, #bbb, #999, #bbb, #999)'
      : isHovered
        ? 'linear-gradient(145deg, #565656, #d5d5d5, #8a8a8a, #cbcbcb, #d5d5d5)'
        : 'linear-gradient(145deg, #787776, #f5f5f5, #aaaaaa, #ebebeb, #f5f5f5)',
    borderStyle: 'solid',
    borderWidth: '5px',
    borderColor: disabled ? '#999' : '#cccccc',
    borderRadius: '25px',
    fontFamily: 'tm-default',
    fontSize: 'clamp(24px, 5vw, 40px)',
    color: disabled ? '#666' : (isHovered ? '#eee' : 'black'),
    textTransform: 'uppercase',
    textAlign: 'center',
    cursor: disabled ? 'default' : 'pointer',
    transition: 'all 0.1s ease',
    opacity: disabled ? 0.6 : 1
  };

  const wrapperStyle = {
    height: '75px',
    width: '100%',
    maxWidth: '500px',
    display: 'block',
    marginLeft: 'auto',
    marginRight: 'auto',
    marginBottom: '30px',
    textAlign: 'center',
    padding: '0 10px'
  };

  return (
    <div style={wrapperStyle}>
      <button
        style={buttonStyle}
        onMouseEnter={() => !disabled && setIsHovered(true)}
        onMouseLeave={() => !disabled && setIsHovered(false)}
        onClick={disabled ? undefined : onClick}
        disabled={disabled}
      >
        {children}
      </button>
    </div>
  );
}

export default LinkButton;