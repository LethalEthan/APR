function generateRandomNumber() {

    var rand = Math.floor(Math.random() * 1000) + 1;
 
    document.getElementById('display').innerText = rand;
  }