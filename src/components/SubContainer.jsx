import React from 'react';
import { formStyles } from '../styles/formStyles';

function SubContainer({ children }) {
  return (
    <div style={{
      textAlign: 'center',
      fontSize: '22px'
    }}>
      {children}
    </div>
  );
}

function SubContainerElement({ children, label, input }) {
  return (
    <div style={{
      maxWidth: '750px',
      margin: '5px auto',
      minHeight: '35px',
      lineHeight: '35px',
      textAlign: 'left'
    }}>
      {label && <label>{label}</label>}
      {input && (
        <input 
          {...input}
          style={{
            ...formStyles.optionInput,
            ...input.style
          }}
        />
      )}
      {children}
    </div>
  );
}

export { SubContainer, SubContainerElement };