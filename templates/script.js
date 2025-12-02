document.addEventListener('DOMContentLoaded', () => {
    const containerWrapper = document.getElementById('class-cards-container-wrapper');

   let containerHtml = `<!-- Example class cards -->
        <div class="col-md-6 col-lg-4">
            <div class="card class-card">
                <div class="card-body">
                    <h5 class="card-title">English Literature</h5>
                    <p class="text-muted mb-2"><i class="bi bi-person"></i> Teacher: item.teacher</p>
                    <p class="text-muted mb-2"><i class="bi bi-clock"></i> item.day_of_week at item.time</p>
                    <p class="text-muted mb-2"><i class="bi bi-people"></i> 12/20 spots filled</p>
                    <p class="card-text mb-3">Explore classic and modern literature through discussion and analysis.</p>
                    <div class="mb-3">
                        <strong class="d-block mb-2">Session Dates:</strong>
                        <span class="badge bg-light text-dark me-1">Nov 6</span>
                        <span class="badge bg-light text-dark me-1">Nov 13</span>
                        <span class="badge bg-light text-dark me-1">Nov 20</span>
                        <span class="badge bg-light text-dark me-1">Nov 27</span>
                    </div>
                    <button class="btn btn-secondary w-100" disabled>
                        <i class="bi bi-lock"></i> Requires Approval
                    </button>
                    <button class="btn btn-outline-primary w-100 mt-2" onclick="showRequestApprovalModal('English Literature')">
                        <i class="bi bi-hand-thumbs-up"></i> Request Enrollment
                    </button>
                </div>
            </div>
        </div>
        <!-- Example class cards -->
        <div class="col-md-6 col-lg-4">
            <div class="card class-card">
                <div class="card-body">
                    <h5 class="card-title">{item.title}</h5>
                    <p class="text-muted mb-2"><i class="bi bi-person"></i> </p>
                    <p class="text-muted mb-2"><i class="bi bi-clock"></i> Wednesdays at 2:00 PM</p>
                    <p class="text-muted mb-2"><i class="bi bi-people"></i> 12/20 spots filled</p>
                    <p class="card-text mb-3">Explore classic and modern literature through discussion and analysis.</p>
                    <div class="mb-3">
                        <strong class="d-block mb-2">Session Dates:</strong>
                        <span class="badge bg-light text-dark me-1">Nov 6</span>
                        <span class="badge bg-light text-dark me-1">Nov 13</span>
                        <span class="badge bg-light text-dark me-1">Nov 20</span>
                        <span class="badge bg-light text-dark me-1">Nov 27</span>
                    </div>
                    <button class="btn btn-secondary w-100" disabled>
                        <i class="bi bi-lock"></i> Requires Approval
                    </button>
                    <button class="btn btn-outline-primary w-100 mt-2" onclick="showRequestApprovalModal('English Literature')">
                        <i class="bi bi-hand-thumbs-up"></i> Request Enrollment
                    </button>
                </div>
            </div>
        </div>

        <div class="col-md-6 col-lg-4">
            <div class="card class-card">
                <div class="card-body">
                    <h5 class="card-title">Science Lab</h5>
                    <p class="text-muted mb-2"><i class="bi bi-person"></i> Teacher: Dr. Michael Chen</p>
                    <p class="text-muted mb-2"><i class="bi bi-clock"></i> Fridays at 3:30 PM</p>
                    <p class="text-muted mb-2"><i class="bi bi-people"></i> 8/15 spots filled</p>
                    <p class="card-text mb-3">Hands-on experiments in physics, chemistry, and biology.</p>
                    <div class="mb-3">
                        <strong class="d-block mb-2">Session Dates:</strong>
                        <span class="badge bg-light text-dark me-1">Nov 1</span>
                        <span class="badge bg-light text-dark me-1">Nov 8</span>
                        <span class="badge bg-light text-dark me-1">Nov 15</span>
                        <span class="badge bg-light text-dark me-1">Nov 22</span>
                    </div>
                    <button class="btn btn-secondary w-100" disabled>
                        <i class="bi bi-lock"></i> Requires Approval
                    </button>
                    <button class="btn btn-outline-primary w-100 mt-2" onclick="showRequestApprovalModal('Science Lab')">
                        <i class="bi bi-hand-thumbs-up"></i> Request Enrollment
                    </button>
                </div>
            </div>
        </div>

        <div class="col-md-6 col-lg-4">
            <div class="card class-card">
                <div class="card-body">
                    <h5 class="card-title">Art Workshop</h5>
                    <p class="text-muted mb-2"><i class="bi bi-person"></i> Teacher: Emily Rodriguez</p>
                    <p class="text-muted mb-2"><i class="bi bi-clock"></i> Tuesdays at 4:00 PM</p>
                    <p class="text-muted mb-2"><i class="bi bi-people"></i> 15/15 spots filled</p>
                    <p class="card-text mb-3">Creative expression through various artistic mediums.</p>
                    <div class="mb-3">
                        <strong class="d-block mb-2">Session Dates:</strong>
                        <span class="badge bg-light text-dark me-1">Nov 5</span>
                        <span class="badge bg-light text-dark me-1">Nov 12</span>
                        <span class="badge bg-light text-dark me-1">Nov 19</span>
                        <span class="badge bg-light text-dark me-1">Nov 26</span>
                    </div>
                    <button class="btn btn-secondary w-100" disabled>Class Full</button>
                </div>
            </div>
        </div>

        <!-- Example: Class student CAN enroll in (same as enrolled class) -->
        <div class="col-md-6 col-lg-4">
            <div class="card class-card border-success">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h5 class="card-title mb-0">Math 101 - December</h5>
                        <span class="badge bg-success">Eligible</span>
                    </div>
                    <p class="text-muted mb-2"><i class="bi bi-person"></i> Teacher: Alex Johnson</p>
                    <p class="text-muted mb-2"><i class="bi bi-clock"></i> Mondays at 10:00 AM</p>
                    <p class="text-muted mb-2"><i class="bi bi-people"></i> 5/20 spots filled</p>
                    <p class="card-text small mb-3">Continue your mathematics journey in December.</p>
                    <div class="mb-3">
                        <strong class="d-block mb-2">Session Dates:</strong>
                        <span class="badge bg-light text-dark me-1">Dec 2</span>
                        <span class="badge bg-light text-dark me-1">Dec 9</span>
                        <span class="badge bg-light text-dark me-1">Dec 16</span>
                        <span class="badge bg-light text-dark me-1">Dec 23</span>
                    </div>
                    <button class="btn btn-primary w-100" onclick="showEnrollModal('Math 101 - December')">
                        <i class="bi bi-check-circle"></i> Enroll Now
                    </button>
                </div>
            </div>
        </div>`

    // Append the new HTML to the wrapper
    containerWrapper.insertAdjacentHTML('beforeend', containerHtml);


    // const apiEndpoint = '/api/your-data-endpoint'; // Replace with your actual backend endpoint

    // fetch(apiEndpoint)
    //     .then(response => response.json())
    //     .then(data => {
    //         // Iterate over the data array and create HTML for each item
    //         data.forEach(item => {
    //             const containerHtml = `
    //                 <div class="col-md-6 col-lg-4">
    //                     <div class="p-4 border bg-light h-100">
    //                         <h2 class="h5">${item.title}</h2>
    //                         <p>${item.content}</p>
    //                         <a href="#" class="btn btn-primary">Learn More</a>
    //                     </div>
    //                 </div>
    //             `;
    //             // Append the new HTML to the wrapper
    //             containerWrapper.insertAdjacentHTML('beforeend', containerHtml);
    //         });
    //     })
    //     .catch(error => {
    //         console.error('Error fetching data:', error);
    //         containerWrapper.innerHTML = '<p class="text-danger">Failed to load content.</p>';
    //     });
});