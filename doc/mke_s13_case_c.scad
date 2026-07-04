// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Case
// Parametric OpenSCAD  |  v6.1 (Case-Only View Modes)
// Units: mm
// =====================================================

include <mke_s13_config_c.scad>
use <threads.scad>

// --- VIEW CONFIGURATION ---
// 1 = Print layout   (bottom shell + print-ready lid, side by side, bed-flat)
// 2 = Assembled      (bottom shell + lid closed, with transparency, CASE ONLY)
// 3 = Cutaway        (sliced along Y-midplane to verify clearances, CASE ONLY)
// 4 = Exploded       (bottom shell + lid separated vertically, CASE ONLY)
view_mode = 1;

$fn = 48;

// =====================================================
// 2D PROFILES
// =====================================================
module pcb_2d() {
    hull() {
        translate([0,        pcb_w/2])         circle(r=0.1);
        translate([chev_l,   0])                circle(r=0.1);
        translate([chev_l,   pcb_w])           circle(r=0.1);
        // X-axis extended to ensure rapid-fab PCBs don't crash into the wall
        translate([pcb_l - corner_r + fab_x_tol, corner_r])           circle(r=corner_r);
        translate([pcb_l - corner_r + fab_x_tol, pcb_w - corner_r])   circle(r=corner_r);
    }
}

module pcb_box_section_2d() {
    intersection() {
        pcb_2d();
        translate([safe_line_x - boolean_overlap, -5])
            square([box_l + 6, pcb_w + 10]);
    }
}

module shell_2d(w) {
    offset(r = w + pcb_gap_y) pcb_box_section_2d();
}

// =====================================================
// BOTTOM SHELL (Upstanding Tongue for Labyrinth - NO BEAD)
// =====================================================
module bottom_shell() {
    difference() {
        union() {
            // 1. Main Outer body box (Stops early to create the outer shoulder)
            linear_extrude(outer_h - lip_h) shell_2d(wall);

            // 2. Upstanding Centered Tongue (Clean, straight wall)
            translate([0, 0, outer_h - lip_h])
                linear_extrude(lip_h)
                    difference() {
                        shell_2d(tongue_out);
                        shell_2d(tongue_in);
                    }
        }

        // 3. Interior hollow (Leaves the internal cavity perfectly clear)
        translate([0, 0, floor_t])
            linear_extrude(outer_h + 1)
                offset(r = pcb_gap_y) pcb_box_section_2d();

        // --- DROP-IN SLOT FOR PCB (flat, rectangular -- slot_gap clearance
        // on all 4 sides, referenced directly off z_pcb_seat so it's exactly
        // slot_gap below the PCB's actual seated bottom surface, not a
        // reconstruction that silently drops solder_z_tol) ---
        translate([
            safe_line_x - boolean_overlap - pcb_gap_y - wall - 1.0,
            -slot_gap,
            z_pcb_seat - slot_gap
        ])
            cube([
                wall + 2.0,
                pcb_w + 2*slot_gap,
                outer_h + 1
            ]);
    }

    // PCB mounting standoffs
    standoff_h = z_pcb_seat - floor_t;
    translate([0, 0, floor_t]) {
        for (y_off = [-hole_sp/2, hole_sp/2]) {
            translate([hole_x, hole_cy + y_off, 0])
            difference() {
                cylinder(d = pcb_boss_d, h = standoff_h);
                translate([0, 0, -thread_cut_undercut])
                    render()
                    metric_thread(
                        diameter = m3_thread_d,
                        pitch = m3_thread_pitch,
                        length = standoff_h + 1,
                        internal = true
                    );
            }
        }
    }

    // --- CONNECTOR / PCB PARTITION BULKHEAD (LOWER HALF) ---
    intersection() {
        translate([partition_x1, -1000, 0])
            cube([partition_t, 2000, z_pcb_seat]);
        linear_extrude(z_pcb_seat)
            offset(r = pcb_gap_y) pcb_box_section_2d();
    }
}
// Interior hollow profile, with the two pillar locations left solid
// so the lid's screw pillars stay continuous from the mating plane
// down through the hollowed cavity to the shoulder near the PCB.
pillar_hollow_margin = 0.4; // small diametral safety margin around
                             // the boss so its edge doesn't sit exactly
                             // coincident with the hollow-cut boundary
module lid_center_hollow_2d() {
    difference() {
        shell_2d(0);
        for (y_off = [-hole_sp/2, hole_sp/2])
            translate([hole_x, hole_cy + y_off])
                circle(d = pcb_boss_d + pillar_hollow_margin);
    }
}
// =====================================================
// LID (Double-Wall Groove Receiver - NO BEAD)
// =====================================================
module lid() {
    // 1. CREATE MAIN SHELL & HOLLOW IT OUT
    difference() {
        union() {
            // Main lid plate
            linear_extrude(lid_t)
                shell_2d(wall + lip_clear + outer_skirt_t);

            // Outer rim extending downwards
            translate([0, 0, -lip_h])
                linear_extrude(lip_h)
                    shell_2d(wall + lip_clear + outer_skirt_t);

            // --- INTEGRATED PILLAR SHOULDERS ---
            for (y_off = [-hole_sp/2, hole_sp/2]) {
                translate([hole_x, hole_cy + y_off, 0]) {
                     shoulder_h = outer_h - (z_pcb_seat + pcb_t) - shoulder_pcb_clearance_z;
                    translate([0, 0, -shoulder_h])
                        cylinder(d = pcb_boss_d + pillar_hollow_margin, h = shoulder_h + lid_reconnect_h);
                }
            }
        }

        // Hollow out the very center
        translate([0, 0, -lip_h - lid_center_hollow_overcut])
            linear_extrude(lip_h + 2*lid_center_hollow_overcut)
                lid_center_hollow_2d();

        // Cut the explicit labyrinth groove
        translate([0, 0, -lip_h - 0.01])
            linear_extrude(lip_h + groove_overcut)
                difference() {
                    shell_2d(tongue_out + lip_clear);
                    shell_2d(tongue_in - lip_clear);
                }

        // TIGHT CONNECTOR SLOT
        translate([
            partition_x2,
            pcb_w/2 - (conn_l / 2) - cable_clear,
            -lip_h - 0.1
         ])
            cube([
                (pcb_l + cable_clear) - partition_x2,
                conn_l + 2*cable_clear,
                lid_t + lip_h + 0.2
            ]);
        // --- M3 COUNTERSUNK CLEARANCE HOLES + PHASE-LOCKED LID PILLAR THREAD ---
        // Each hole is now 4 zones stacked tip-to-plate: [plain clearance]
        // -> [threaded section] -> [plain clearance] -> [countersink cone].
        // Previously the whole shaft was plain 3.4mm clearance, so only
        // the bottom standoff's thread (~2.3mm engagement) gripped the
        // screw -- thin for M3 in FDM plastic. Threading part of the lid
        // pillar too adds real engagement length, but ONLY works if it's
        // in phase with the bottom standoff's thread:
        //
        // A physical screw is one continuous helix -- its thread crest is
        // a fixed function of Z (crest angle advances 360 deg every
        // m3_thread_pitch mm) that doesn't reset partway down the shank.
        // Two separate female thread cuts mesh with the SAME screw at
        // once, without cross-threading, only if the distance between
        // their local z=0 reference planes is an exact whole number of
        // pitches -- matching diameter/pitch alone isn't enough. The
        // bottom standoff's thread (bottom_shell(), above) starts its
        // local z=0 at global z = floor_t - thread_cut_undercut. The lid
        // pillar thread below is placed at the nearest global Z at or
        // past the tip's lead-in buffer that is an exact multiple of
        // m3_thread_pitch away from that same reference -- computed, not
        // hand-picked, so it stays correct if PCB/clearance dimensions in
        // the config ever change.
        for (y_off = [-hole_sp/2, hole_sp/2]) {
            translate([hole_x, hole_cy + y_off, 0]) {
                shoulder_h = outer_h - (z_pcb_seat + pcb_t) - shoulder_pcb_clearance_z;

                // Phase-aligned lid-pillar thread placement (all global Z,
                // then converted to this pillar's local frame, where z=0
                // is the lid's underside mating plane = global z=outer_h)
                bottom_thread_origin_global = floor_t - thread_cut_undercut;
                pillar_tip_global       = z_pcb_seat + pcb_t + shoulder_pcb_clearance_z; // == outer_h - shoulder_h
                lid_thread_lead_in      = 0.6; // plain clearance left at the fragile pillar tip before threading starts
                lid_thread_top_clear    = 1.5; // plain clearance left just under the lid plate, above the thread
                raw_thread_start        = pillar_tip_global + lid_thread_lead_in;
                n_pitches_up            = ceil((raw_thread_start - bottom_thread_origin_global) / m3_thread_pitch);
                lid_thread_start_global = bottom_thread_origin_global + n_pitches_up * m3_thread_pitch;
                lid_thread_len          = (outer_h - lid_thread_top_clear) - lid_thread_start_global;
                lid_thread_start_local  = lid_thread_start_global - outer_h;

                // 1. Lower plain clearance shaft: pillar tip up to thread start
                translate([0, 0, -shoulder_h - 0.1])
                    cylinder(d = 3.4, h = lid_thread_start_local - (-shoulder_h - 0.1) + thread_cut_undercut, $fn = 50);

                // 2. Phase-locked internal M3 thread
                translate([0, 0, lid_thread_start_local])
                    render()
                    metric_thread(
                        diameter = m3_thread_d,
                        pitch = m3_thread_pitch,
                        length = lid_thread_len,
                        internal = true
                    );

                // 3. Upper plain clearance shaft: thread top up through the rest of the pillar and lid plate
                translate([0, 0, lid_thread_start_local + lid_thread_len - thread_cut_undercut])
                    cylinder(d = 3.4, h = (lid_t + 0.1) - (lid_thread_start_local + lid_thread_len) + (2 * thread_cut_undercut), $fn = 50);

                // 4. The countersink cone for the flat head screw
                // Top surface of the lid is at Z = lid_t
                // d1 is the bottom of the cone (3.4mm), d2 is the top (6.2mm)
                translate([0, 0, lid_t - 1.7])
                    cylinder(h = 1.7 + 0.1, d1 = 3.4, d2 = 6.2, $fn = 50);
            }
        }

    } // <-- THE DIFFERENCE BLOCK ENDS HERE

    // 2. ADD INTERNAL FEATURES (Now safe from being hollowed out)

    // --- CLOSURE BAFFLE (flat, rectangular) ---
    // baffle_h is the exact drop from the lid mating plane (outer_h) down to
    // slot_gap above the PCB's actual top surface (z_pcb_seat + pcb_t):
    //   outer_h - (z_pcb_seat + pcb_t + slot_gap)
    //     = (inner_h + floor_t) - (floor_t + pcb_t + pcb_z_gap + solder_z_tol + pcb_t + slot_gap)
    //     = inner_h - (2*pcb_t + pcb_z_gap + slot_gap + solder_z_tol)
    baffle_h = inner_h - (2 * pcb_t + pcb_z_gap + slot_gap + solder_z_tol);
    // Narrower than the shell's slot by clearance_per_side/side, purely so the
    // baffle drops into the shell's slot without binding on its side walls --
    // unrelated to PCB clearance, which is handled by baffle_h alone.
    baffle_w = pcb_w + 2*slot_gap - (2 * clearance_per_side);

    translate([
        safe_line_x - boolean_overlap - pcb_gap_y - wall,
        pcb_w/2 - baffle_w/2,
        -baffle_h
    ])
        cube([
            wall,
            baffle_w,
            baffle_h + lid_reconnect_h + groove_overcut
        ]);

    // --- CONNECTOR / PCB PARTITION BULKHEAD (UPPER HALF) ---
    partition_drop_h = outer_h - (z_pcb_seat + pcb_t + slot_gap);
    partition_w      = pcb_w + 2*pcb_gap_y - baffle_clearance;

    translate([
        partition_x1,
        pcb_w/2 - partition_w/2,
        -partition_drop_h
    ])
        cube([partition_t, partition_w, partition_drop_h + lid_reconnect_h]);
}

// =====================================================
// PRINT-READY ORIENTATIONS
// =====================================================
module lid_print_ready() {
    translate([0, 0, lid_t])
        mirror([0, 0, 1])
            lid();
}

// =====================================================
// VIEW-MODE RENDERING MODULES (CASE ONLY -- no PCB/connectors)
// =====================================================
module print_layout() {
    plate_gap = 2.0;
    bottom_shell_max_y = pcb_w + 2*(wall + pcb_gap_y);
    lid_min_y_offset    = wall + lip_clear + outer_skirt_t + pcb_gap_y;

    color("SteelBlue",  0.85) bottom_shell();
    color("LightBlue",  0.70)
        translate([0, bottom_shell_max_y + lid_min_y_offset + plate_gap, 0])
            lid_print_ready();
}

module full_system() {
    color("SteelBlue", 0.65) bottom_shell();
    translate([0, 0, outer_h]) color("LightBlue", 0.50) lid();
}

module exploded_system() {
    gap = 25;

    color("SteelBlue", 0.85) bottom_shell();
    translate([0, 0, outer_h + gap]) color("LightBlue", 0.75) lid();
}

// =====================================================
// RENDER EXECUTION
// =====================================================
if (view_mode == 1) {
    print_layout();
} else if (view_mode == 2) {
    full_system();
} else if (view_mode == 3) {
    difference() {
        full_system();
        translate([-10, pcb_w/2, -5]) {
            cube([pcb_l + 30, pcb_w + 20, outer_h + 30]);
        }
    }
} else if (view_mode == 4) {
    exploded_system();
}