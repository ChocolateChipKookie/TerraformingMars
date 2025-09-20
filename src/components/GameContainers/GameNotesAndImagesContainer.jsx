import React, { useRef } from 'react';
import Container from '../Container';
import { SubContainer, SubContainerElement } from '../Container';
import styles from '../../styles/GamePage.module.css';

function GameNotesAndImagesContainer({ 
  gameConfig, 
  dispatch,
  readOnly = false
}) {
  const fileInputRef = useRef(null);

  const handleImageSelect = (e) => {
    const files = Array.from(e.target.files);
    
    files.forEach(file => {
      if (gameConfig.images.length >= 4) {
        alert('Maximum 4 images allowed');
        return;
      }
      
      if (!file.type.startsWith('image/')) {
        alert('Please select only image files');
        return;
      }
      
      const reader = new FileReader();
      reader.onload = (event) => {
        const base64Data = event.target.result.split(',')[1]; // Remove data:image/...;base64, prefix
        dispatch({
          type: 'ADD_IMAGE',
          image: {
            image_data: base64Data,
            mime_type: file.type,
            preview: event.target.result, // Keep full data URL for preview
            name: file.name
          }
        });
      };
      reader.readAsDataURL(file);
    });
    
    // Clear the input so the same file can be selected again
    e.target.value = '';
  };

  const handleRemoveImage = (index) => {
    dispatch({ type: 'REMOVE_IMAGE', index });
  };

  return (
    <Container title="Notes & Images" titleStyle="banner">
      <SubContainer>
        <SubContainerElement>
          <label>Game Notes:</label>
          <textarea
            className={styles.notesTextarea}
            value={gameConfig.note}
            onChange={(e) => dispatch({ type: 'SET_NOTE', value: e.target.value })}
            placeholder="Add any notes about this game..."
            rows={4}
            readOnly={readOnly}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Game Images ({gameConfig.images.length}/4):</label>
          
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            onChange={handleImageSelect}
            style={{ display: 'none' }}
          />

          <div className={styles.imageGrid}>
            {gameConfig.images.map((image, index) => (
                <div key={index} className={styles.imageContainer}>
                  {readOnly ? (
                    // In view mode - make images clickable
                    <a
                      href={`http://localhost:8080/api/images/${image.id}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className={styles.imageLink}
                    >
                      <img
                        src={`http://localhost:8080/api/images/${image.id}`}
                        alt={`Game image ${index + 1}`}
                        className={styles.gameImage}
                        onError={(e) => {
                          // Fallback for broken image links
                          e.target.style.backgroundColor = '#f0f0f0';
                          e.target.style.display = 'flex';
                          e.target.style.alignItems = 'center';
                          e.target.style.justifyContent = 'center';
                          e.target.innerHTML = 'Image not found';
                        }}
                      />
                    </a>
                  ) : (
                    // In edit mode - images not clickable
                    <img
                      src={
                        image.preview
                          ? image.preview // For newly uploaded images (has preview)
                          : `http://localhost:8080/api/images/${image.id}` // For existing images from backend
                      }
                      alt={`Game image ${index + 1}`}
                      className={styles.gameImage}
                      onError={(e) => {
                        // Fallback for broken image links
                        e.target.style.backgroundColor = '#f0f0f0';
                        e.target.style.display = 'flex';
                        e.target.style.alignItems = 'center';
                        e.target.style.justifyContent = 'center';
                        e.target.innerHTML = 'Image not found';
                      }}
                    />
                  )}
                  {!readOnly && (
                    <button
                      type="button"
                      onClick={() => handleRemoveImage(index)}
                      className={styles.removeImageButton}
                    >
                      ×
                    </button>
                  )}
                </div>
              ))}
            {!readOnly && gameConfig.images.length < 4 && (
              <div
                className={styles.addImageBox}
                onClick={() => fileInputRef.current?.click()}
              >
                <span className={styles.addImagePlus}>+</span>
              </div>
            )}
          </div>
        </SubContainerElement>
      </SubContainer>
    </Container>
  );
}

export default GameNotesAndImagesContainer;