const rl = @import("raylib");
const std = @import("std");

var palette = [_]struct { u8, u8, u8, u8 }{
    .{ 0, 0, 0, 0 }, // 0: Clear
    .{ 7, 7, 7, 255 }, // 1: Dark Gray
    .{ 31, 7, 7, 255 }, // 2: Dark Red
    .{ 47, 15, 7, 255 }, // 3: Red-Brown
    .{ 71, 15, 7, 255 }, // 4: Brown
    .{ 87, 23, 7, 255 }, // 5: Orange-Brown
    .{ 103, 31, 7, 255 }, // 6: Dark Orange
    .{ 119, 31, 7, 255 }, // 7: Orange
    .{ 143, 39, 7, 255 }, // 8: Bright Orange
    .{ 159, 47, 7, 255 }, // 9: Light Orange
    .{ 175, 63, 7, 255 }, // 10: Lighter Orange
    .{ 191, 71, 7, 255 }, // 11: Reddish Orange
    .{ 199, 71, 7, 255 }, // 12: Bright Reddish Orange
    .{ 223, 79, 7, 255 }, // 13: Red-Orange
    .{ 223, 87, 7, 255 }, // 14: Fiery Orange
    .{ 223, 87, 7, 255 }, // 15: Fiery Orange (repeat for gradient)
    .{ 215, 95, 7, 255 }, // 16: Yellowish Orange
    .{ 215, 103, 15, 255 }, // 17: Yellow-Orange
    .{ 207, 111, 15, 255 }, // 18: Yellow
    .{ 207, 119, 15, 255 }, // 19: Bright Yellow
    .{ 207, 127, 15, 255 }, // 20: Pale Yellow
    .{ 207, 135, 23, 255 }, // 21: Pale Yellowish White
    .{ 199, 135, 23, 255 }, // 22: Lighter Yellowish White
    .{ 199, 143, 23, 255 }, // 23: Yellow-White
    .{ 199, 151, 31, 255 }, // 24: Pale Yellow-White
    .{ 191, 159, 31, 255 }, // 25: Bright Yellow-White
    .{ 191, 159, 31, 255 }, // 26: Bright Yellow-White (repeat for gradient)
    .{ 191, 167, 39, 255 }, // 27: Very Bright Yellow
    .{ 191, 167, 39, 255 }, // 28: Very Bright Yellow (repeat for gradient)
    .{ 255, 255, 63, 255 }, // 29: Near White
    .{ 255, 255, 111, 255 }, // 30: Faint Yellowish White
    .{ 255, 255, 159, 255 }, // 31: Faint Pale White
    .{ 255, 255, 191, 255 }, // 32: Pale White
    .{ 255, 255, 223, 255 }, // 33: Near Full White
    .{ 255, 255, 239, 255 }, // 34: Almost Full White
    .{ 255, 255, 247, 255 }, // 35: Subtle White
    .{ 255, 255, 255, 255 }, // 36: Full White
};

const Game = struct {
    al: std.mem.Allocator,
    framebuffer: [][]u8 = undefined, // 2D array to store fire intensity values
    screenWidth: u32 = 0,
    screenHeight: u32 = 0,
    fireWidth: u32 = 0,
    fireHeight: u32 = 0,
    texture: rl.RenderTexture2D = undefined,

    pub fn init(al: std.mem.Allocator, width: i32, height: i32) !Game {
        const fireH = @divFloor(height, 12);
        rl.setWindowSize(width, fireH + 50);
        var game = Game{
            .al = al,
            .screenWidth = @intCast(width),
            .screenHeight = @intCast(fireH + 50),
            .fireWidth = @intCast(width),
            .fireHeight = @intCast(fireH),
            .texture = try rl.loadRenderTexture(width, fireH),
        };
        game.resetFramebuffer();
        return game;
    }

    fn free(self: *Game) void {
        for (self.framebuffer) |*row| {
            self.al.free(row.*);
        }
        self.al.free(self.framebuffer);
    }

    fn resetFramebuffer(self: *Game) void {
        self.framebuffer = self.al.alloc([]u8, self.fireHeight) catch unreachable;
        for (self.framebuffer) |*row| {
            row.* = self.al.alloc(u8, self.fireWidth) catch unreachable;
            // Clear the buffer
            for (0..self.fireWidth) |x| {
                row.*[x] = 0;
            }
        }

        // Set the bottom row of the framebuffer to the maximum intensity
        for (0..self.fireWidth) |x| {
            self.framebuffer[self.fireHeight - 1][x] = palette.len - 1;
        }
    }

    fn spreadFire(self: *Game, from: usize) void {
        const rand_offset: i8 = std.crypto.random.intRangeAtMost(i8, 0, 3) - 1;
        const delta: i64 = @as(i64, @intCast(from)) - self.fireWidth;
        const to: i64 = delta + @as(i64, rand_offset);

        // Validate `to` is within bounds
        if (to < 0 or to >= @as(i64, self.fireWidth * self.fireHeight)) return;

        // Calculate source and destination coordinates
        const src_x = from % self.fireWidth;
        const src_y = from / self.fireWidth;

        const dst_x = @mod(to, @as(i64, self.fireWidth));
        const dst_y = @divFloor(to, @as(i64, self.fireWidth));

        if (dst_x < 0 or dst_x >= @as(i64, self.fireWidth)) return;

        const old_value = self.framebuffer[src_y][src_x];
        const rand_delta = std.crypto.random.intRangeAtMost(u8, 0, 2);

        const new_value = if (old_value > rand_delta) old_value - rand_delta else 0;
        self.framebuffer[@as(usize, @intCast(dst_y))][@as(usize, @intCast(dst_x))] = new_value;
    }

    fn doFire(self: *Game) void {
        for (0..self.fireWidth) |x| {
            for (0..self.fireHeight) |y| {
                self.spreadFire(y * self.fireWidth + x);
            }
        }
    }

    pub fn update(self: *Game) void {
        self.doFire();
    }

    pub fn draw(self: *Game) void {
        const pixel_data = self.al.alloc(u8, self.fireWidth * self.fireHeight * 4) catch unreachable;
        defer self.al.free(pixel_data);

        // Prepare pixel data based on the framebuffer
        for (0..self.fireHeight) |y| {
            for (0..self.fireWidth) |x| {
                const intensity = self.framebuffer[y][x];
                const idx = (y * self.fireWidth + x) * 4;

                if (intensity == 0) {
                    // Transparent black for areas without flame
                    pixel_data[idx + 0] = 0; // Red
                    pixel_data[idx + 1] = 0; // Green
                    pixel_data[idx + 2] = 0; // Blue
                    pixel_data[idx + 3] = 0; // Alpha (fully transparent)
                } else {
                    // Flame color
                    pixel_data[idx + 0] = palette[intensity][0]; // Red
                    pixel_data[idx + 1] = palette[intensity][1]; // Green
                    pixel_data[idx + 2] = palette[intensity][2]; // Blue
                    pixel_data[idx + 3] = 255; // Fully opaque
                }
            }
        }

        // Update the render texture with the pixel data
        rl.updateTexture(self.texture.texture, pixel_data.ptr);

        // Draw the render texture's texture to the screen
        rl.drawTexturePro(
            self.texture.texture,
            rl.Rectangle{
                .x = 0,
                .y = 0,
                .width = @floatFromInt(self.fireWidth),
                .height = @floatFromInt(@as(i32, @intCast(self.fireHeight))),
            },
            rl.Rectangle{
                .x = 0,
                .y = 0,
                .width = @floatFromInt(self.screenWidth),
                .height = @floatFromInt(self.screenHeight),
            },
            rl.Vector2{ .x = 0, .y = 0 },
            0.0,
            rl.Color.white,
        );
    }
};

pub fn main() !void {
    const initial_width = 650;
    const initial_height = 800;

    rl.setTraceLogLevel(.err);

    rl.setConfigFlags(.{ .window_undecorated = true, .window_transparent = true });

    rl.initWindow(initial_width, initial_height, "Doom Fire - Raylib Zig");
    defer rl.closeWindow();

    const h = rl.getMonitorHeight(0);
    const w = rl.getMonitorWidth(0);

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer if (gpa.deinit() == .leak) {
        std.process.exit(1);
    };
    var game = try Game.init(allocator, @intCast(w), @intCast(h));
    defer game.free();

    rl.setWindowPosition(0, h - @as(i32, @intCast(game.screenHeight)));
    rl.setTargetFPS(60);

    while (!rl.windowShouldClose()) {
        const start_time = rl.getTime();

        // Update game logic
        game.update();

        // Draw game
        rl.beginDrawing();
        defer rl.endDrawing();

        rl.clearBackground(rl.Color.blank);

        game.draw();

        const end_time = rl.getTime();
        const frame_time = (end_time - start_time) * 1000.0; // Convert to milliseconds

        const fps = if (frame_time > 0) 1000.0 / frame_time else 0.0;

        // Print frame time
        std.debug.print("Frame Time: {d:.2} ms | FPS: {d:.2}\r", .{ frame_time, fps });
    }
    std.debug.print("\n", .{});
}
