import React from 'react';
import { formStyles } from '../styles/formStyles';

function SubContainer({ children }) {
  return (
    <div style={{
      textAlign: 'center',
      fontSize: '2.5rem',
      padding: '0.5rem 2rem'
    }}>
      {children}
    </div>
  );
}

function SubContainerElement({ children, label, input }) {
  return (
    <div style={{
      maxWidth: '750px',
      margin: '0.5rem auto',
      minHeight: '3rem',
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
