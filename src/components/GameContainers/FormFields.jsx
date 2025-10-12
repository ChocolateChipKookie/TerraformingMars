import React from 'react';
import styles from '../../styles/GamePage.module.css';

// SelectField - handles dropdowns with read-only support
export function SelectField({ 
  value, 
  options, 
  onChange, 
  readOnly = false, 
  placeholder = "Select...",
  disabled = false,
  className = styles.containerInput,
  displayFormatter = null,
  children
}) {
  if (readOnly) {
    // Format the display value
    let displayValue = value;
    const isNotAchieved = value === "" || value == null || value === -1;

    if (displayFormatter) {
      displayValue = displayFormatter(value, options);
    } else if (isNotAchieved) {
      displayValue = placeholder;
    } else if (children) {
      // If using children for options, just display the value
      displayValue = value;
    } else {
      // Find the option label if options is an array of objects
      const option = options?.find?.(opt =>
        typeof opt === 'object' ? opt.value === value : opt === value
      );
      displayValue = option ? (typeof option === 'object' ? option.label : option) : value;
    }

    return (
      <div
        className={className}
        style={{
          cursor: 'default',
          userSelect: 'none',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          textAlign: 'center',
          opacity: isNotAchieved ? 0.5 : 1,
          color: isNotAchieved ? '#888' : 'inherit'
        }}
      >
        {displayValue}
      </div>
    );
  }
  
  return (
    <select
      className={className}
      value={value ?? ""}
      onChange={onChange}
      disabled={disabled}
    >
      {(value === "" || value == null) && placeholder && <option value="">{placeholder}</option>}
      {children || options?.map((option) => {
        const optionValue = typeof option === 'object' ? option.value : option;
        const optionLabel = typeof option === 'object' ? option.label : option;
        return (
          <option key={optionValue} value={optionValue}>
            {optionLabel}
          </option>
        );
      })}
    </select>
  );
}

// InputField - handles text/number inputs with read-only support
export function InputField({ 
  value, 
  onChange, 
  readOnly = false,
  type = "text",
  placeholder = "",
  disabled = false,
  className = styles.containerInput,
  inputMode,
  pattern,
  required,
  style
}) {
  if (readOnly) {
    return (
      <div 
        className={className} 
        style={{ 
          cursor: 'default', 
          userSelect: 'none',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          textAlign: 'center',
          ...style
        }}
      >
        {value || placeholder}
      </div>
    );
  }
  
  return (
    <input
      type={type}
      className={className}
      value={value}
      onChange={onChange}
      disabled={disabled}
      placeholder={placeholder}
      inputMode={inputMode}
      pattern={pattern}
      required={required}
      style={style}
    />
  );
}