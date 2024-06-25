document.addEventListener('DOMContentLoaded',function(){
    fetch('/Properties').then(response=>response.json()).then(data=>{
        document.getElementById('message').textContent=data.message;
    })
})