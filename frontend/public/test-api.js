// Simple test to verify API connectivity
console.log('Testing API connectivity...');

const testAPI = async () => {
  try {
    console.log('Making request to http://localhost:8080/api/v1/health');
    const response = await fetch('http://localhost:8080/api/v1/health', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      mode: 'cors',
    });
    
    console.log('Response status:', response.status);
    console.log('Response headers:', response.headers);
    
    if (response.ok) {
      const data = await response.json();
      console.log('✅ API Response:', data);
    } else {
      console.error('❌ API Error:', response.status, response.statusText);
    }
  } catch (error) {
    console.error('❌ Fetch Error:', error);
  }
};

testAPI();
