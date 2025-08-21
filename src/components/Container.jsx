import React from 'react';

function Container({ children, title, titleStyle = "page-title" }) {
  return (
    <div style={{
      background: 'linear-gradient(145deg, #787776, #f5f5f5, #aaaaaa, #ebebeb, #f5f5f5)',
      borderStyle: 'solid',
      borderWidth: '5px',
      borderColor: '#cccccc',
      borderRadius: '20px',
      maxWidth: '900px',
      fontFamily: 'tm-default',
      marginLeft: 'auto',
      marginRight: 'auto',
      marginTop: '1%',
      marginBottom: '1%',
      textAlign: 'center',
      textTransform: 'uppercase',
      color: 'black'
    }}>
      {title && titleStyle === "page-title" && (
        <div style={{
          fontSize: '4rem',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center'
        }}>
          {title}
        </div>
      )}
      {title && titleStyle === "banner" && (
        <div style={{
          marginTop: '10px',
          borderTop: '5px black solid',
          borderBottom: '5px black solid',
          fontSize: '3rem',
          backgroundColor: 'rgb(240, 240, 240)'
        }}>
          {title}
        </div>
      )}
      {children}
    </div>
  );
}

export default Container;
