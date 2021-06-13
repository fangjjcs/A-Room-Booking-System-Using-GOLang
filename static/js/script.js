
var forms = document.querySelectorAll(".needs-validation");

(function () {
  "use strict";

  // Fetch all the forms we want to apply custom Bootstrap validation styles to
  
  // Loop over them and prevent submission
  Array.prototype.slice.call(forms).forEach(function (form) {
    form.addEventListener(
      "submit",
      function (event) {
        if (!form.checkValidity()) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add("was-validated");
      },
      false
    );
  });
})();

const elem = document.getElementById("reservation-dates");
const rangepicker = new DateRangePicker(elem, {
  format: "yyyy-mm-dd",
});



let attention = Prompt();

document.getElementById("btnPrompt").addEventListener("click", function(){
  // attention.toast({msg:"success!"});
  // 
  // Ajax : fetching a Json Object
  fetch("/search-availability-json")
    .then(res => res.json())
    .then(data => {
      console.log(data);
    })

})


function Prompt() {
  let toast = function (c) {
      
    const {
        msg= "",
        icon= "success",
        position= "top-end"
    } = c;
      
    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 1000,
      timerProgressBar: false,
      didOpen: (toast) => {
        toast.addEventListener("mouseenter", Swal.stopTimer);
        toast.addEventListener("mouseleave", Swal.resumeTimer);
      },
    });

    Toast.fire({});
  };
    
  return {
      toast: toast,
  }
}


