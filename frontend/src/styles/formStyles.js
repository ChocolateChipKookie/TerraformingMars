// Common form styles for the Terraforming Mars app
export const formStyles = {
  optionInput: {
    textAlign: 'center',
    fontFamily: 'inherit',
    fontSize: '1.8rem',
    background: 'inherit',
    float: 'right',
    height: '3rem',
    width: '50%',
    boxSizing: 'border-box'
  },

  subcontainerBox: {
    borderStyle: 'solid',
    borderWidth: '5px',
    borderColor: 'black',
    borderRadius: '25px',
    marginTop: '20px',
    marginBottom: '20px',
    backgroundColor: 'rgb(240, 240, 240)',
    display: 'block'
  },

  checkboxLabel: {
    display: 'block',
    lineHeight: '3rem',
    height: '3rem',
    width: '300px',
    position: 'relative',
    fontSize: '22px',
    userSelect: 'none',
    WebkitUserSelect: 'none',
    MozUserSelect: 'none',
    msUserSelect: 'none'
  },

  hiddenCheckbox: {
    position: 'absolute',
    opacity: 0,
    cursor: 'pointer',
    height: 0,
    width: 0
  },

  expansionIconStyle: {
    width: "45px",
    height: "45px",
    borderRadius: "50%",
    border: "2px solid #000",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    fontSize: "28px",
    fontWeight: "bold",
    userSelect: "none",
    flexShrink: 0,
  },

  checkmark: {
    position: 'absolute',
    top: 0,
    left: '275px',
    width: '25px',
    height: '25px'
  },

  playerInputDiv: {
    textAlign: 'center',
    verticalAlign: 'center',
    display: 'flex',
    justifyContent: 'space-between',
    margin: '0.3rem 0'
  },

  containerInput: {
    borderColor: '#666',
    borderRadius: '4px',
    textAlign: 'center',
    fontWeight: 'bolder',
    fontFamily: 'inherit',
    fontSize: '1.8rem',
    backgroundColor: 'white',
    width: '48%',
    height: '3rem',
    boxSizing: 'border-box',
    display: 'inline-block'
  },

  milestoneLabel: {
    borderColor: '#666',
    textAlign: 'center',
    fontWeight: 'bolder',
    fontFamily: 'inherit',
    fontSize: '1.8rem',
    width: '48%',
    height: '3rem',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  }
};
